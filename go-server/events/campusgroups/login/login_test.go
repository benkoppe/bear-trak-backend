package login

import (
	"log"
	"os"
	"testing"
)

func TestLogin(t *testing.T) {
	url := "https://cornell.campusgroups.com/login_only"
	p := LoginParams{
		LoginEmail:       os.Getenv("LOGIN_EMAIL"),
		OtpEmail:         os.Getenv("OTP_EMAIL"),
		OtpEmailPassword: os.Getenv("OTP_EMAIL_PASSWORD"),
	}

	out, err := GetLoginCookie(url, p)
	if err != nil {
		t.Fatalf("login failed: %v \n", err)
	}

	log.Printf("got: %s \n", out)

	out2, err := GetLoginCookie(url, p)
	if err != nil {
		t.Fatalf("login failed: %v \n", err)
	}

	log.Printf("second run (test caching), got: %s \n", out2)
}
