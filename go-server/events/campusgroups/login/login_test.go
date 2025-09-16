package login

import (
	"log"
	"os"
	"testing"
)

func TestLogin(t *testing.T) {
	loginURL := "https://cornell.campusgroups.com/login_only"
	loginEmail := os.Getenv("LOGIN_EMAIL")
	otpEmail := os.Getenv("OTP_EMAIL")
	otpEmailPassword := os.Getenv("OTP_EMAIL_PASSWORD")

	out, err := GetLoginCookie(loginURL, loginEmail, otpEmail, otpEmailPassword)
	if err != nil {
		t.Fatalf("login failed: %v \n", err)
	}

	log.Printf("got: %s \n", out)

	out2, err := GetLoginCookie(loginURL, loginEmail, otpEmail, otpEmailPassword)
	if err != nil {
		t.Fatalf("login failed: %v \n", err)
	}

	log.Printf("second run (test caching), got: %s \n", out2)
}
