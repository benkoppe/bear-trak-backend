package login

import (
	"fmt"
	"regexp"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

// Regex to extract the 6-digit code from the subject
var otpRegex = regexp.MustCompile(`Sign In Code \((\d{6})\)`)

func fetchRecentEnvelopes(c *client.Client, n uint32) (map[uint32]*imap.Envelope, error) {
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		return nil, fmt.Errorf("failed to select INBOX: %w", err)
	}
	if mbox.Messages == 0 {
		return map[uint32]*imap.Envelope{}, nil
	}

	from := uint32(1)
	if mbox.Messages > n {
		from = mbox.Messages - n + 1
	}
	seqset := new(imap.SeqSet)
	seqset.AddRange(from, mbox.Messages)

	messages := make(chan *imap.Message, n)
	if err := c.Fetch(seqset, []imap.FetchItem{imap.FetchUid, imap.FetchEnvelope}, messages); err != nil {
		return nil, fmt.Errorf("fetch failed: %w", err)
	}

	out := make(map[uint32]*imap.Envelope)
	for msg := range messages {
		if msg.Envelope != nil {
			out[msg.Uid] = msg.Envelope
		}
	}
	return out, nil
}

func waitForNewOtp(c *client.Client, baseline map[uint32]*imap.Envelope, wait time.Duration) (string, error) {
	deadline := time.Now().Add(wait)
	for time.Now().Before(deadline) {
		current, err := fetchRecentEnvelopes(c, 10)
		if err != nil {
			return "", err
		}

		for uid, env := range current {
			if _, existed := baseline[uid]; !existed && env != nil {
				if match := otpRegex.FindStringSubmatch(env.Subject); len(match) == 2 {
					otp := match[1]

					// --- Mark this message as read ---
					seqset := new(imap.SeqSet)
					seqset.AddNum(uid)
					item := imap.FormatFlagsOp(imap.AddFlags, true)
					flags := []any{imap.SeenFlag}
					if err := c.UidStore(seqset, item, flags, nil); err != nil {
						return "", fmt.Errorf("failed to mark OTP email as seen: %w", err)
					}

					return otp, nil
				}
			}
		}

		time.Sleep(2 * time.Second)
	}
	return "", fmt.Errorf("timed out waiting for OTP email")
}
