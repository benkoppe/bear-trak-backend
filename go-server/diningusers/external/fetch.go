// Package external loads all external CBORD dining data.
package external

import (
	"encoding/json"
	"fmt"

	"github.com/benkoppe/bear-trak-backend/go-server/utils"
)

// FetchUserID - nil if authentication fails, errors for other failure
func FetchUserID(baseURL string, sessionID string) (*UserIDResponseBody, error) {
	fullURL, err := utils.ExtendURL(baseURL, "user")
	if fullURL == nil {
		return nil, fmt.Errorf("failed to extend url: %w", err)
	}

	requestBody := map[string]any{
		"method": "retrieve",
		"params": map[string]string{
			"sessionId": sessionID,
		},
	}

	resp, err := utils.DoPostRequest[userIDResponse](*fullURL, requestBody)
	if resp == nil {
		return nil, err
	}

	return resp.Response, nil
}

func FetchBarcodeSeed(baseURL string, sessionID string, institutionID string) (*string, error) {
	fullURL, err := utils.ExtendURL(baseURL, "configuration")
	if fullURL == nil {
		return nil, fmt.Errorf("failed to extend url: %w", err)
	}

	requestBody := map[string]any{
		"method": "nativeStartup",
		"params": map[string]string{
			"clientType":    "ios",
			"clientVersion": "4.33.23",
			"institutionId": institutionID,
			"sessionId":     sessionID,
		},
	}

	resp, err := utils.DoPostRequest[barcodeSeedResponse](*fullURL, requestBody)
	if resp == nil {
		return nil, err
	}

	return &resp.Response.BarcodeSeed, err
}

func FetchCashlessKey(baseURL string, sessionID string) (*string, error) {
	fullURL, err := utils.ExtendURL(baseURL, "user")
	if fullURL == nil {
		return nil, fmt.Errorf("failed to extend url: %w", err)
	}

	requestBody := map[string]any{
		"method": "retrieveSetting",
		"params": map[string]string{
			"settingName": "CashlessKey",
			"sessionId":   sessionID,
		},
	}

	resp, err := utils.DoPostRequest[cashlessKeyResponse](*fullURL, requestBody)
	if resp == nil {
		return nil, err
	}

	return &resp.Response.Value, err
}

func FetchAccounts(baseURL string, sessionID string, userID string) (*accountsResponseBody, error) {
	fullURL, err := utils.ExtendURL(baseURL, "commerce")
	if fullURL == nil {
		return nil, fmt.Errorf("failed to extend url: %w", err)
	}

	requestBody := map[string]any{
		"method": "retrieveAccountsByUser",
		"params": map[string]string{
			"sessionId": sessionID,
			"userId":    userID,
		},
	}

	resp, err := utils.DoPostRequest[accountsResponse](*fullURL, requestBody)
	if resp == nil {
		return nil, err
	}

	return resp.Response, err
}

func FetchUserPhoto(baseURL string, sessionID string, userID string) (*userPhotoResponseBody, error) {
	fullURL, err := utils.ExtendURL(baseURL, "user")
	if fullURL == nil {
		return nil, fmt.Errorf("failed to extend url: %w", err)
	}

	requestBody := map[string]any{
		"method": "retrieveUserPhoto",
		"params": map[string]string{
			"sessionId": sessionID,
			"userId":    userID,
		},
	}

	resp, err := utils.DoPostRequest[userPhotoResponse](*fullURL, requestBody)
	if resp == nil {
		return nil, err
	}

	return resp.Response, nil
}

func FetchDisplayTenders(baseURL string, sessionID string, institutionID string) ([]string, error) {
	fullURL, err := utils.ExtendURL(baseURL, "configuration")
	if fullURL == nil {
		return nil, fmt.Errorf("failed to extend url: %w", err)
	}

	requestBody := map[string]any{
		"method": "retrieveSetting",
		"params": map[string]string{
			"domain":        "get",
			"category":      "feature",
			"name":          "display_tenders",
			"institutionId": institutionID,
			"sessionId":     sessionID,
		},
	}

	resp, err := utils.DoPostRequest[retrieveSettingResponse](*fullURL, requestBody)
	if resp == nil {
		return nil, err
	}

	// decode tender IDs from json
	var tenderIDs []string
	err = json.Unmarshal([]byte(resp.Response.Value), &tenderIDs)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling tender IDs: %w", err)
	}

	return tenderIDs, nil
}
