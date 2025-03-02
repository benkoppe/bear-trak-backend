package dining_users

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
	"unicode"

	"github.com/benkoppe/bear-trak-backend/go-server/api"
	"github.com/benkoppe/bear-trak-backend/go-server/db"
	"github.com/benkoppe/bear-trak-backend/go-server/dining-users/external"
	"golang.org/x/sync/errgroup"
)

func hashUserId(userIdResp external.UserIDResponseBody) string {
	hasher := sha256.New()
	hasher.Write([]byte(userIdResp.ID))
	return hex.EncodeToString(hasher.Sum(nil))
}

func CreateUser(ctx context.Context, externalBaseUrl string, params api.PostV1DiningUserParams, queries *db.Queries) (api.PostV1DiningUserRes, error) {
	idResp, err := external.FetchUserID(externalBaseUrl, params.SessionId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user info: %w", err)
	}
	if idResp == nil {
		return &api.PostV1DiningUserUnauthorized{}, nil
	}

	hashedId := hashUserId(*idResp)
	users, err := queries.GetDiningUserAll(ctx, hashedId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users from database: %w", err)
	}
	if len(users) > 0 {
		// disallow more than one login per user
		return &api.PostV1DiningUserBadRequest{}, nil
	}

	resp, err := external.CreatePIN(externalBaseUrl, params.DeviceId, params.PIN, params.SessionId)
	if err != nil {
		return nil, fmt.Errorf("failed to create PIN: %w", err)
	}

	if resp == nil {
		return &api.PostV1DiningUserUnauthorized{}, nil
	}

	_, err = queries.CreateDiningUser(ctx, db.CreateDiningUserParams{
		UserID:   hashedId,
		DeviceID: params.DeviceId,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user in database: %w", err)
	}

	user, err := GetUser(externalBaseUrl, api.GetV1DiningUserParams{SessionId: params.SessionId})
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if diningUser, ok := user.(*api.DiningUser); ok {
		resp := api.PostV1DiningUserCreated{
			Message:   "User created.",
			ID:        diningUser.ID,
			FirstName: diningUser.FirstName,
			LastName:  diningUser.LastName,
			PhotoJpeg: diningUser.PhotoJpeg,
		}
		return &resp, nil
	}

	return nil, fmt.Errorf("failed to convert user to DiningUser")
}

func DeleteUser(ctx context.Context, externalBaseUrl string, params api.DeleteV1DiningUserParams, queries *db.Queries) (api.DeleteV1DiningUserRes, error) {
	idResp, err := external.FetchUserID(externalBaseUrl, params.SessionId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user info: %w", err)
	}
	if idResp == nil {
		return &api.DeleteV1DiningUserUnauthorized{}, nil
	}

	hashedId := hashUserId(*idResp)
	err = queries.DeleteDiningUser(ctx, hashedId)
	if err != nil {
		return nil, fmt.Errorf("failed to delete user from database: %w", err)
	}

	return &api.Success{Message: "User deleted."}, nil
}

func RefreshUserToken(ctx context.Context, externalBaseUrl string, params api.GetV1DiningUserSessionParams, queries *db.Queries) (api.GetV1DiningUserSessionRes, error) {
	user, err := queries.GetDiningUser(ctx, params.DeviceId)
	if err != nil {
		return &api.GetV1DiningUserSessionUnauthorized{}, nil
	}

	resp, err := external.CreateSession(externalBaseUrl, params.DeviceId, params.PIN)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	if resp == nil {
		return &api.GetV1DiningUserSessionUnauthorized{}, nil
	}

	queries.UpdateDiningUserSession(ctx, user.ID)

	res := api.GetV1DiningUserSessionOKApplicationJSON(*resp)
	return &res, nil
}

func GetUser(externalBaseUrl string, params api.GetV1DiningUserParams) (api.GetV1DiningUserRes, error) {
	idResp, err := external.FetchUserID(externalBaseUrl, params.SessionId)
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
		photoResp, err := external.FetchUserPhoto(externalBaseUrl, params.SessionId, idResp.ID)
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
		barcodeSeed, err := external.FetchBarcodeSeed(externalBaseUrl, params.SessionId)
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
		cashlessKey, err := external.FetchCashlessKey(externalBaseUrl, params.SessionId)
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

	// extract account type (first word) and filter
	var response []api.DiningUserAccount
	for _, account := range resp.Accounts {
		accountType, shortName := splitAccountName(account)
		account.Name = shortName

		if strings.HasPrefix(accountType, "CC1") || strings.HasPrefix(accountType, "GET") || strings.HasPrefix(accountType, "01n") {
			continue
		}

		if strings.HasPrefix(accountType, "CB") || strings.HasPrefix(accountType, "BRB") {
			response = append(response, convertExternalAccount(account, true))
			continue
		}

		response = append(response, convertExternalAccount(account, false))
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
