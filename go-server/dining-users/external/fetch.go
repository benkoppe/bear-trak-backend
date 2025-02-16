package external

import (
	"fmt"

	"github.com/benkoppe/bear-trak-backend/go-server/utils"
)

// nil if authentication fails, errors for other failure
func FetchUserID(baseUrl string, sessionId string) (*userIDResponseBody, error) {
	fullUrl, err := utils.ExtendUrl(baseUrl, "user")
	if fullUrl == nil {
		return nil, fmt.Errorf("failed to extend url: %w", err)
	}

	requestBody := map[string]interface{}{
		"method": "retrieve",
		"params": map[string]string{
			"sessionId": sessionId,
		},
	}

	resp, err := utils.DoPostRequest[userIDResponse](*fullUrl, requestBody)
	if resp == nil {
		return nil, err
	}

	return resp.Response, nil
}

func FetchBarcode(baseUrl string, sessionId string) (*string, error) {
	fullUrl, err := utils.ExtendUrl(baseUrl, "authentication")
	if fullUrl == nil {
		return nil, fmt.Errorf("failed to extend url: %w", err)
	}

	requestBody := map[string]interface{}{
		"method": "retrievePatronBarcodePayload",
		"params": map[string]string{
			"sessionId": sessionId,
		},
	}

	resp, err := utils.DoPostRequest[stringResponse](*fullUrl, requestBody)
	if resp == nil {
		return nil, err
	}

	return resp.Response, err
}

func FetchAccounts(baseUrl string, sessionId string, userId string) (*accountsResponseBody, error) {
	fullUrl, err := utils.ExtendUrl(baseUrl, "commerce")
	if fullUrl == nil {
		return nil, fmt.Errorf("failed to extend url: %w", err)
	}

	requestBody := map[string]interface{}{
		"method": "retrieveAccountsByUser",
		"params": map[string]string{
			"sessionId": sessionId,
			"userId":    userId,
		},
	}

	resp, err := utils.DoPostRequest[accountsResponse](*fullUrl, requestBody)
	if resp == nil {
		return nil, err
	}

	return resp.Response, err
}
