package external

import (
	"fmt"

	"github.com/benkoppe/bear-trak-backend/go-server/utils"
)

func DeletePIN(baseUrl string, deviceId string, sessionId string) (*bool, error) {
	fullUrl, err := utils.ExtendUrl(baseUrl, "user")
	if fullUrl == nil {
		return nil, fmt.Errorf("failed to extend url: %w", err)
	}

	requestBody := map[string]interface{}{
		"method": "deletePIN",
		"params": map[string]string{
			"sessionId": sessionId,
			"deviceId":  deviceId,
		},
	}

	resp, err := utils.DoPostRequest[boolResponse](*fullUrl, requestBody)
	if resp == nil {
		return nil, err
	}

	return resp.Response, nil
}
