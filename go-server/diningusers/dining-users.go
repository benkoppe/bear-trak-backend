// Package diningusers includes all methods for CBORD dining integration.
package diningusers

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"unicode"

	"github.com/benkoppe/bear-trak-backend/go-server/api"
	"github.com/benkoppe/bear-trak-backend/go-server/db"
	"github.com/benkoppe/bear-trak-backend/go-server/diningusers/external"
	"github.com/benkoppe/bear-trak-backend/go-server/utils"
	"golang.org/x/sync/errgroup"
)

func hashUserID(userIDResp external.UserIDResponseBody) string {
	hasher := sha256.New()
	hasher.Write([]byte(userIDResp.ID))
	return hex.EncodeToString(hasher.Sum(nil))
}

func CreateUser(ctx context.Context, externalBaseURL, institutionID string, params api.PostV1DiningUserParams, queries *db.Queries) (api.PostV1DiningUserRes, error) {
	idResp, err := external.FetchUserID(externalBaseURL, params.SessionId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user info: %w", err)
	}
	if idResp == nil {
		return &api.PostV1DiningUserUnauthorized{}, nil
	}

	hashedID := hashUserID(*idResp)
	users, err := queries.GetDiningUserAll(ctx, hashedID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users from database: %w", err)
	}
	if len(users) > 0 {
		// disallow more than one login per user
		return &api.PostV1DiningUserBadRequest{}, nil
	}

	resp, err := external.CreatePIN(externalBaseURL, params.DeviceId, params.PIN, params.SessionId)
	if err != nil {
		return nil, fmt.Errorf("failed to create PIN: %w", err)
	}

	if resp == nil {
		return &api.PostV1DiningUserUnauthorized{}, nil
	}

	_, err = queries.CreateDiningUser(ctx, db.CreateDiningUserParams{
		UserID:   hashedID,
		DeviceID: params.DeviceId,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user in database: %w", err)
	}

	user, err := GetUser(externalBaseURL, institutionID, api.GetV1DiningUserParams{SessionId: params.SessionId})
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if diningUser, ok := user.(*api.DiningUser); ok {
		resp := api.PostV1DiningUserCreated{
			Message:        "User created.",
			ID:             diningUser.ID,
			FirstName:      diningUser.FirstName,
			LastName:       diningUser.LastName,
			PhotoJpeg:      diningUser.PhotoJpeg,
			BarcodeSeedHex: diningUser.BarcodeSeedHex,
			CashlessKey:    diningUser.CashlessKey,
		}
		return &resp, nil
	}

	return nil, fmt.Errorf("failed to convert user to DiningUser")
}

func DeleteUser(ctx context.Context, externalBaseURL string, params api.DeleteV1DiningUserParams, queries *db.Queries) (api.DeleteV1DiningUserRes, error) {
	idResp, err := external.FetchUserID(externalBaseURL, params.SessionId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user info: %w", err)
	}
	if idResp == nil {
		return &api.DeleteV1DiningUserUnauthorized{}, nil
	}

	hashedID := hashUserID(*idResp)
	err = queries.DeleteDiningUser(ctx, hashedID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete user from database: %w", err)
	}

	return &api.Success{Message: "User deleted."}, nil
}

func RefreshUserToken(ctx context.Context, externalBaseURL string, params api.GetV1DiningUserSessionParams, queries *db.Queries) (api.GetV1DiningUserSessionRes, error) {
	user, err := queries.GetDiningUser(ctx, params.DeviceId)
	if err != nil {
		return &api.GetV1DiningUserSessionUnauthorized{}, nil
	}

	resp, err := external.CreateSession(externalBaseURL, params.DeviceId, params.PIN)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	if resp == nil {
		return &api.GetV1DiningUserSessionUnauthorized{}, nil
	}

	err = queries.UpdateDiningUserSession(ctx, user.ID)
	if err != nil {
		log.Printf("failed to update dining user session timestamp: %v", err)
	}

	res := api.GetV1DiningUserSessionOKApplicationJSON(*resp)
	return &res, nil
}

func GetUser(externalBaseURL, institutionID string, params api.GetV1DiningUserParams) (api.GetV1DiningUserRes, error) {
	idResp, err := external.FetchUserID(externalBaseURL, params.SessionId)
	if err != nil {
		return nil, fmt.Errorf("failed to get user id: %w", err)
	}
	if idResp == nil {
		return &api.GetV1DiningUserUnauthorized{}, nil
	}

	res := convertExternalUser(*idResp)

	// user errgroup to handle concurrent operations
	var eg errgroup.Group

	// fetch photo concurrently
	eg.Go(func() error {
		photoResp, err := external.FetchUserPhoto(externalBaseURL, params.SessionId, idResp.ID)
		if photoResp != nil {
			if photoResp.MimeType == "image/jpeg" && photoResp.Data != "" {
				decodedBytes, err := base64.StdEncoding.DecodeString(photoResp.Data)
				if err == nil {
					res.PhotoJpeg = decodedBytes
				}
			}
		} else {
			if err != nil {
				fmt.Println("photoResp had non-breaking error: %w", err)
			} else {
				fmt.Println("photoResp is nil")
			}
		}
		return nil // photo errors are non-breaking
	})

	// fetch barcode seed concurrently
	eg.Go(func() error {
		barcodeSeed, err := external.FetchBarcodeSeed(externalBaseURL, params.SessionId, institutionID)
		if err != nil {
			return fmt.Errorf("failed to get barcode seed; %w", err)
		}
		if barcodeSeed != nil {
			res.BarcodeSeedHex = *barcodeSeed
		}
		return nil
	})

	// fetch cashless key concurrently
	eg.Go(func() error {
		cashlessKey, err := external.FetchCashlessKey(externalBaseURL, params.SessionId)
		if err != nil {
			return fmt.Errorf("failed to get cashless key; %w", err)
		}
		if cashlessKey != nil {
			res.CashlessKey = *cashlessKey
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return &res, nil
}

func convertExternalUser(user external.UserIDResponseBody) api.DiningUser {
	return api.DiningUser{
		ID:        user.ID,
		FirstName: toProperCaseEachWord(user.FirstName),
		LastName:  toProperCaseEachWord(user.LastName),
	}
}

func toProperCaseWord(word string) string {
	if word == "" {
		return word
	}

	runes := []rune(strings.ToLower(word))
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func toProperCaseEachWord(s string) string {
	words := strings.Fields(s)
	for i, w := range words {
		words[i] = toProperCaseWord(w)
	}
	return strings.Join(words, " ")
}

func GetUserBarcode(externalBaseURL string, params api.GetV1DiningUserBarcodeParams) (api.GetV1DiningUserBarcodeRes, error) {
	resp, err := external.FetchBarcode(externalBaseURL, params.SessionId)
	if err != nil {
		return nil, fmt.Errorf("failed to get user barcode: %w", err)
	}

	if resp == nil {
		return &api.GetV1DiningUserBarcodeUnauthorized{}, nil
	}

	res := api.GetV1DiningUserBarcodeOKApplicationJSON(*resp)
	return &res, nil
}

func GetUserAccounts(externalBaseURL, institutionID string, params api.GetV1DiningUserAccountsParams) (api.GetV1DiningUserAccountsRes, error) {
	idResp, err := external.FetchUserID(externalBaseURL, params.SessionId)
	if err != nil {
		return nil, fmt.Errorf("failed to get user id: %w", err)
	}
	if idResp == nil {
		return &api.GetV1DiningUserAccountsUnauthorized{}, nil
	}

	resp, err := external.FetchAccounts(externalBaseURL, params.SessionId, idResp.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user accounts: %w", err)
	}

	if resp == nil {
		return &api.GetV1DiningUserAccountsUnauthorized{}, nil
	}

	displayTenders, err := external.FetchDisplayTenders(externalBaseURL, params.SessionId, institutionID)
	if err != nil {
		fmt.Printf("failed to get display tenders: %v\n", err)
	}

	// extract account type (first word) and filter
	var response []api.DiningUserAccount
	for _, account := range resp.Accounts {
		_, shortName := splitAccountName(account)
		account.Name = shortName

		// we should only display account tenders in displayTenders
		if !utils.Contains(displayTenders, account.Tender) {
			continue
		}

		moneyBalance := true
		if account.Type == 1 {
			moneyBalance = false
		}

		response = append(response, convertExternalAccount(account, moneyBalance))
	}

	res := api.GetV1DiningUserAccountsOKApplicationJSON(response)
	return &res, nil
}

func convertExternalAccount(account external.Account, moneyBalance bool) api.DiningUserAccount {
	var balance api.DiningUserAccountBalance
	if moneyBalance {
		balance = api.NewDiningUserMoneyBalanceDiningUserAccountBalance(api.DiningUserMoneyBalance{
			Money: account.Balance,
		})
	} else {
		balance = api.NewDiningUserSwipeBalanceDiningUserAccountBalance(api.DiningUserSwipeBalance{
			Swipes: int(account.Balance),
		})
	}

	return api.DiningUserAccount{
		AccountId: account.ID,
		Name:      account.Name,
		Balance:   balance,
	}
}

func splitAccountName(account external.Account) (firstWord, remaining string) {
	parts := strings.SplitN(account.Name, " ", 2)
	firstWord = parts[0]
	if len(parts) > 1 {
		remaining = parts[1]
	}
	return
}
