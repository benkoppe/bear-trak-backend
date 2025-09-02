package external

import (
	"fmt"

	"github.com/benkoppe/bear-trak-backend/go-server/utils"
)

func DeletePIN(baseURL string, deviceID string, sessionID string) (*bool, error) {
	fullURL, err := utils.ExtendURL(baseURL, "user")
	if fullURL == nil {
		return nil, fmt.Errorf("failed to extend url: %w", err)
	}

	requestBody := map[string]any{
		"method": "deletePIN",
		"params": map[string]string{
			"sessionId": sessionID,
			"deviceId":  deviceID,
		},
	}

	resp, err := utils.DoPostRequest[boolResponse](*fullURL, requestBody)
	if resp == nil {
		return nil, err
	}

	return resp.Response, nil
}
