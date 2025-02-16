package external

import (
	"fmt"

	"github.com/benkoppe/bear-trak-backend/go-server/utils"
)

func CreatePIN(baseUrl string, deviceId string, pin string, sessionId string) (*bool, error) {
	fullUrl, err := utils.ExtendUrl(baseUrl, "user")
	if fullUrl == nil {
		return nil, fmt.Errorf("failed to extend url: %w", err)
	}

	requestBody := map[string]interface{}{
		"method": "createPIN",
		"params": map[string]string{
			"sessionId": sessionId,
			"deviceId":  deviceId,
			"PIN":       pin,
		},
	}

	resp, err := utils.DoPostRequest[boolResponse](*fullUrl, requestBody)
	if resp == nil {
		return nil, err
	}

	return resp.Response, nil
}

func CreateSession(baseUrl string, deviceId string, pin string) (*string, error) {
	fullUrl, err := utils.ExtendUrl(baseUrl, "authentication")
	if fullUrl == nil {
		return nil, fmt.Errorf("failed to extend url: %w", err)
	}

	requestBody := map[string]interface{}{
		"method": "authenticatePIN",
		"params": map[string]interface{}{
			"systemCredentials": map[string]string{
				"domain":   "",
				"userName": "get_mobile",
				"password": "NOTUSED",
			},
			"deviceId": deviceId,
			"pin":      pin,
		},
	}

	resp, err := utils.DoPostRequest[stringResponse](*fullUrl, requestBody)
	if resp == nil {
		return nil, err
	}

	return resp.Response, nil
}
