package external

import (
	"fmt"

	"github.com/benkoppe/bear-trak-backend/go-server/utils"
)

func CreatePIN(baseURL string, deviceID string, pin string, sessionID string) (*bool, error) {
	fullURL, err := utils.ExtendURL(baseURL, "user")
	if fullURL == nil {
		return nil, fmt.Errorf("failed to extend url: %w", err)
	}

	requestBody := map[string]any{
		"method": "createPIN",
		"params": map[string]string{
			"sessionId": sessionID,
			"deviceId":  deviceID,
			"PIN":       pin,
		},
	}

	resp, err := utils.DoPostRequest[boolResponse](*fullURL, requestBody)
	if resp == nil {
		return nil, err
	}

	return resp.Response, nil
}

func CreateSession(baseURL string, deviceID string, pin string) (*string, error) {
	fullURL, err := utils.ExtendURL(baseURL, "authentication")
	if fullURL == nil {
		return nil, fmt.Errorf("failed to extend url: %w", err)
	}

	requestBody := map[string]any{
		"method": "authenticatePIN",
		"params": map[string]any{
			"systemCredentials": map[string]string{
				"domain":   "",
				"userName": "get_mobile",
				"password": "NOTUSED",
			},
			"deviceId": deviceID,
			"pin":      pin,
		},
	}

	resp, err := utils.DoPostRequest[stringResponse](*fullURL, requestBody)
	if resp == nil {
		return nil, err
	}

	return resp.Response, nil
}
