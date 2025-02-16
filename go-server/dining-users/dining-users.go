package dining_users

import (
	"fmt"
	"strings"

	"github.com/benkoppe/bear-trak-backend/go-server/api"
	"github.com/benkoppe/bear-trak-backend/go-server/dining-users/external"
)

func CreateUser(externalBaseUrl string, params api.PostV1DiningUserParams) (api.PostV1DiningUserRes, error) {
	resp, err := external.CreatePIN(externalBaseUrl, params.DeviceId, params.PIN, params.SessionId)
	if err != nil {
		return nil, fmt.Errorf("failed to create PIN: %w", err)
	}

	if resp == nil {
		return &api.PostV1DiningUserUnauthorized{}, nil
	}

	return &api.Success{Message: "User created."}, nil
}

func DeleteUser(externalBaseUrl string, session api.DeleteV1DiningUserParams) (api.DeleteV1DiningUserRes, error) {
	// TODO: implement correctly
	// this needs a database to do correctly, to associate user ids with devices

	return &api.Success{Message: "User deleted."}, nil
}

func RefreshUserToken(externalBaseUrl string, params api.GetV1DiningUserSessionParams) (api.GetV1DiningUserSessionRes, error) {
	resp, err := external.CreateSession(externalBaseUrl, params.DeviceId, params.PIN)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	if resp == nil {
		return &api.GetV1DiningUserSessionUnauthorized{}, nil
	}

	res := api.GetV1DiningUserSessionOKApplicationJSON(*resp)
	return &res, nil
}

func GetUserBarcode(externalBaseUrl string, params api.GetV1DiningUserBarcodeParams) (api.GetV1DiningUserBarcodeRes, error) {
	resp, err := external.FetchBarcode(externalBaseUrl, params.SessionId)
	if err != nil {
		return nil, fmt.Errorf("failed to get user barcode: %w", err)
	}

	if resp == nil {
		return &api.GetV1DiningUserBarcodeUnauthorized{}, nil
	}

	res := api.GetV1DiningUserBarcodeOKApplicationJSON(*resp)
	return &res, nil
}

func GetUserAccounts(externalBaseUrl string, params api.GetV1DiningUserAccountsParams) (api.GetV1DiningUserAccountsRes, error) {
	idResp, err := external.FetchUserID(externalBaseUrl, params.SessionId)
	if err != nil {
		return nil, fmt.Errorf("failed to get user id: %w", err)
	}
	if idResp == nil {
		return &api.GetV1DiningUserAccountsUnauthorized{}, nil
	}

	resp, err := external.FetchAccounts(externalBaseUrl, params.SessionId, idResp.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user accounts: %w", err)
	}

	if resp == nil {
		return &api.GetV1DiningUserAccountsUnauthorized{}, nil
	}

	var accounts []api.DiningUserAccount
	for _, account := range resp.Accounts {
		accounts = append(accounts, convertExternalAccount(account))
	}

	// extract account type (first word) and filter
	var filteredAccounts []api.DiningUserAccount
	for _, account := range accounts {
		accountType, shortName := splitAccountName(account)
		account.Name = shortName

		if strings.HasPrefix(accountType, "CB") || strings.HasPrefix(accountType, "BRB") || strings.HasPrefix(accountType, "GET") {
			filteredAccounts = append(filteredAccounts, account)
		}
	}

	res := api.GetV1DiningUserAccountsOKApplicationJSON(filteredAccounts)
	return &res, nil
}

func convertExternalAccount(account external.Account) api.DiningUserAccount {
	return api.DiningUserAccount{
		AccountId: account.ID,
		Name:      account.Name,
		Balance:   account.Balance,
	}
}

func splitAccountName(account api.DiningUserAccount) (firstWord, remaining string) {
	parts := strings.SplitN(account.Name, " ", 2)
	firstWord = parts[0]
	if len(parts) > 1 {
		remaining = parts[1]
	}
	return
}
