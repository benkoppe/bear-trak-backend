// Package external loads all external CBORD dining data.
package external

import (
	"encoding/json"
	"fmt"

	"github.com/benkoppe/bear-trak-backend/go-server/utils"
)

// nil if authentication fails, errors for other failure
func FetchUserID(baseUrl string, sessionId string) (*UserIDResponseBody, error) {
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

func FetchBarcodeSeed(baseUrl string, sessionId string, institutionId string) (*string, error) {
	fullUrl, err := utils.ExtendUrl(baseUrl, "configuration")
	if fullUrl == nil {
		return nil, fmt.Errorf("failed to extend url: %w", err)
	}

	requestBody := map[string]interface{}{
		"method": "nativeStartup",
		"params": map[string]string{
			"clientType":    "ios",
			"clientVersion": "4.33.23",
			"institutionId": institutionId,
			"sessionId":     sessionId,
		},
	}

	resp, err := utils.DoPostRequest[barcodeSeedResponse](*fullUrl, requestBody)
	if resp == nil {
		return nil, err
	}

	return &resp.Response.BarcodeSeed, err
}

func FetchCashlessKey(baseUrl string, sessionId string) (*string, error) {
	fullUrl, err := utils.ExtendUrl(baseUrl, "user")
	if fullUrl == nil {
		return nil, fmt.Errorf("failed to extend url: %w", err)
	}

	requestBody := map[string]interface{}{
		"method": "retrieveSetting",
		"params": map[string]string{
			"settingName": "CashlessKey",
			"sessionId":   sessionId,
		},
	}

	resp, err := utils.DoPostRequest[cashlessKeyResponse](*fullUrl, requestBody)
	if resp == nil {
		return nil, err
	}

	return &resp.Response.Value, err
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

func FetchUserPhoto(baseUrl string, sessionId string, userId string) (*userPhotoResponseBody, error) {
	fullUrl, err := utils.ExtendUrl(baseUrl, "user")
	if fullUrl == nil {
		return nil, fmt.Errorf("failed to extend url: %w", err)
	}

	requestBody := map[string]interface{}{
		"method": "retrieveUserPhoto",
		"params": map[string]string{
			"sessionId": sessionId,
			"userId":    userId,
		},
	}

	resp, err := utils.DoPostRequest[userPhotoResponse](*fullUrl, requestBody)
	if resp == nil {
		return nil, err
	}

	return resp.Response, nil
}

func FetchDisplayTenders(baseUrl string, sessionId string, institutionId string) ([]string, error) {
	fullUrl, err := utils.ExtendUrl(baseUrl, "configuration")
	if fullUrl == nil {
		return nil, fmt.Errorf("failed to extend url: %w", err)
	}

	requestBody := map[string]interface{}{
		"method": "retrieveSetting",
		"params": map[string]string{
			"domain":        "get",
			"category":      "feature",
			"name":          "display_tenders",
			"institutionId": institutionId,
			"sessionId":     sessionId,
		},
	}

	resp, err := utils.DoPostRequest[retrieveSettingResponse](*fullUrl, requestBody)
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
