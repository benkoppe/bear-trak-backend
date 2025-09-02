package email

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"sort"
	"sync"
	"time"

	"github.com/benkoppe/bear-trak-backend/go-server/utils"
	"github.com/benkoppe/bear-trak-backend/go-server/utils/timeutils"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/jhillyerd/enmime"
	openrouter "github.com/revrost/go-openrouter"
	"github.com/revrost/go-openrouter/jsonschema"
)

type MealItem struct {
	ItemName string `json:"itemName"`
}
type MealCategory struct {
	CategoryName string     `json:"categoryName"`
	Items        []MealItem `json:"items"`
}
type Result struct {
	DinnerName string         `json:"dinnerName"`
	MenuFound  bool           `json:"menuFound"`
	Categories []MealCategory `json:"categories"`
}
type DatedResult struct {
	Wednesday time.Time
	Subject   string
	Menu      Result
}

type Cache = *utils.Cache[[]DatedResult]

func InitCache(emailPassword, mistralAPIKey, openrouterAPIKey, openrouterModel string) Cache {
	return utils.NewCache("diningEmail",
		24*time.Hour,
		func() ([]DatedResult, error) {
			return fetchData(emailPassword, mistralAPIKey, openrouterAPIKey, openrouterModel)
		})
}

func menuCacheKey(subject string, date time.Time) string {
	return fmt.Sprintf("%s:%s", subject, date.Format("2006-01-02"))
}

var (
	globalMenuCache = make(map[string]*Result)
	cacheMu         sync.RWMutex
)

func fetchData(emailPassword, mistralAPIKey, openrouterAPIKey, openrouterModel string) ([]DatedResult, error) {
	c, err := client.DialTLS("imap.gmail.com:993", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to IMAP server: %w", err)
	}
	defer c.Logout()

	if err := c.Login("beartrakhousedinners@gmail.com", emailPassword); err != nil {
		return nil, fmt.Errorf("failed to login: %w", err)
	}

	subjects := []string{
		"keeton",
		"cook",
		"becker",
		"rose",
		"bethe",
	}

	var dinners []DatedResult

	for _, subject := range subjects {
		menu, err := getLatestMenu(c, subject, mistralAPIKey, openrouterAPIKey, openrouterModel)
		if err != nil {
			log.Println("Error fetching menu for", subject, ":", err)
			continue
		}
		dinners = append(dinners, *menu)
	}

	return dinners, nil
}

func getLatestMenu(c *client.Client, subject string, mistralAPIKey, openrouterAPIKey, openrouterModel string) (*DatedResult, error) {
	email, err := emailsContainingSubject(c, subject)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch emails: %w", err)
	}

	if len(email) == 0 {
		return nil, fmt.Errorf("no emails")
	}

	latest := email[0]
	nextWednesday := findNextWednesday(latest.InternalDate)

	cacheKey := menuCacheKey(subject, nextWednesday)
	cacheMu.RLock()
	cachedMenu, ok := globalMenuCache[cacheKey]
	cacheMu.RUnlock()
	if ok {
		return &DatedResult{
			Menu:      *cachedMenu,
			Wednesday: nextWednesday,
			Subject:   subject,
		}, nil
	}
	log.Println("Fetching new menu for", subject, "on", nextWednesday)

	html, err := extractHTML(latest)
	if err != nil {
		return nil, fmt.Errorf("failed to extract HTML: %w", err)
	}

	pdf, err := htmlToPDF(html)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to PDF: %w", err)
	}

	uploadURL, err := uploadPDF(pdf, subject+".pdf")
	if err != nil {
		return nil, fmt.Errorf("failed to upload PDF: %w", err)
	}

	ocr, err := sendOCRRequest(uploadURL, mistralAPIKey)
	if err != nil {
		return nil, fmt.Errorf("failed to perform OCR: %w", err)
	}

	menu, err := sendAIRequest(ocr, openrouterAPIKey, openrouterModel)
	if err != nil {
		return nil, fmt.Errorf("failed to send AI request: %w", err)
	}

	cacheMu.Lock()
	globalMenuCache[cacheKey] = menu
	cacheMu.Unlock()

	return &DatedResult{
		Menu:      *menu,
		Wednesday: nextWednesday,
		Subject:   subject,
	}, nil
}

func emailsContainingSubject(c *client.Client, subject string) ([]*imap.Message, error) {
	_, err := c.Select(subject, false)
	if err != nil {
		return nil, fmt.Errorf("failed to select inbox: %w", err)
	}

	criteria := imap.NewSearchCriteria()
	criteria.Since = time.Now().AddDate(0, 0, -30)

	uids, err := c.UidSearch(criteria)
	if err != nil {
		return nil, fmt.Errorf("failed to search for emails: %w", err)
	}
	if len(uids) == 0 {
		return nil, nil
	}

	seqset := new(imap.SeqSet)
	seqset.AddNum(uids...)

	section := &imap.BodySectionName{}

	messages := make(chan *imap.Message, len(uids))
	done := make(chan error, 1)
	go func() {
		done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope, imap.FetchInternalDate, imap.FetchUid, section.FetchItem()}, messages)
	}()

	var results []*imap.Message
	for msg := range messages {
		results = append(results, msg)
	}

	if err := <-done; err != nil {
		return nil, fmt.Errorf("failed to fetch emails: %w", err)
	}

	// sort by date
	sort.Slice(results, func(i, j int) bool {
		return results[i].InternalDate.After(results[j].InternalDate)
	})

	return results, nil
}

func findNextWednesday(startDate time.Time) time.Time {
	currDate := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())

	currentWeekday := currDate.Weekday()
	targetWeekday := time.Wednesday

	daysToAdd := 0
	if currentWeekday < targetWeekday {
		daysToAdd = int(targetWeekday - currentWeekday)
	} else {
		daysToAdd = int(7 - currentWeekday + targetWeekday)
	}

	nextWednesday := currDate.AddDate(0, 0, daysToAdd)

	est := timeutils.LoadEST()
	nextWednesdayEST := time.Date(
		nextWednesday.Year(), nextWednesday.Month(), nextWednesday.Day(), 18, 30, 0, 0, est,
	)
	return nextWednesdayEST
}

func extractHTML(msg *imap.Message) (string, error) {
	r := msg.GetBody(&imap.BodySectionName{})
	if r == nil {
		return "", fmt.Errorf("no body section")
	}

	env, err := enmime.ReadEnvelope(r)
	if err != nil {
		return "", err
	}

	html := env.HTML
	if html == "" {
		html = "<pre>" + env.Text + "</pre>"
	}

	return html, nil
}

func htmlToPDF(html string) ([]byte, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var pdfData []byte
	err := chromedp.Run(ctx,
		chromedp.Navigate("data:text/html,"+url.PathEscape(html)),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			pdfData, _, err = page.PrintToPDF().WithPrintBackground(true).Do(ctx)
			return err
		}),
	)
	if err != nil {
		return nil, err
	}

	return pdfData, nil
}

func uploadPDF(pdfData []byte, fileName string) (string, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	if err := writer.WriteField("reqtype", "fileupload"); err != nil {
		return "", fmt.Errorf("failed to write reqtype field: %w", err)
	}
	if err := writer.WriteField("time", "1h"); err != nil {
		return "", fmt.Errorf("failed to write time field: %w", err)
	}

	// Create file form field.
	fileField, err := writer.CreateFormFile("fileToUpload", fileName)
	if err != nil {
		return "", fmt.Errorf("failed to create file field: %w", err)
	}

	if _, err := fileField.Write(pdfData); err != nil {
		return "", fmt.Errorf("failed to write PDF data: %w", err)
	}

	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer: %w", err)
	}

	req, err := http.NewRequest("POST", "https://litterbox.catbox.moe/resources/internals/api.php", &buf)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	return string(respBody), nil
}

func sendOCRRequest(pdfURL string, mistralAPIKey string) (string, error) {
	payload := map[string]any{
		"model": "mistral-ocr-latest",
		"document": map[string]string{
			"document_url": pdfURL,
			"type":         "document_url",
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal OCR request body: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.mistral.ai/v1/ocr", bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("failed to create OCR request: %w", err)
	}

	// Replace with actual header key if different
	req.Header.Set("Authorization", "Bearer "+mistralAPIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute OCR request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read OCR response: %w", err)
	}

	respStr := string(respBody)
	if respStr == "" {
		return "", fmt.Errorf("empty OCR response")
	}
	return respStr, nil
}

func sendAIRequest(menuOCR string, openrouterAPIKey, openrouterModel string) (*Result, error) {
	ctx := context.Background()

	client := openrouter.NewClient(
		openrouterAPIKey,
	)

	var result Result
	schema, err := jsonschema.GenerateSchemaForType(result)
	if err != nil {
		return nil, fmt.Errorf("GenerateSchemaForType error: %v", err)
	}

	request := openrouter.ChatCompletionRequest{
		Model: openrouterModel,
		Messages: []openrouter.ChatCompletionMessage{
			{
				Role: openrouter.ChatMessageRoleUser,
				Content: openrouter.Content{Text: `What is the entirety of this week's
					house dinner? Use the OCR below to instruct your answer and
					transcribe to the given JSON schema. Don't include special characters
					like â€, the TM symbol, or HTML tags like <br>. Don't include the
					text 'House Dinner' in the dinner name. Sometimes, newsletters don't
					contain house dinner menus. Look carefully to ensure that a house
					dinner menu actually exists. House dinners are only special dinners
					on Wednesday nights, and a menu will look like a complete list of
					options, so be strict about what counts as a menu. If no menu can be
					found, set MenuFound to FALSE in the schema, otherwise set it to TRUE.\n\n` + menuOCR},
			},
		},
		ResponseFormat: &openrouter.ChatCompletionResponseFormat{
			Type: openrouter.ChatCompletionResponseFormatTypeJSONSchema,
			JSONSchema: &openrouter.ChatCompletionResponseFormatJSONSchema{
				Name:   "house-menu",
				Schema: schema,
				Strict: true,
			},
		},
	}

	res, err := client.CreateChatCompletion(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("CreateChatCompletion error: %v", err)
	}

	choices := res.Choices
	if len(choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	choice := choices[0]
	text := choice.Message.Content.Text

	err = json.Unmarshal([]byte(text), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response JSON: %v", err)
	}

	return &result, nil
}
