package external

type boolResponse struct {
	Response *bool `json:"response"`
}

type stringResponse struct {
	Response *string `json:"response"`
}

type userIDResponse struct {
	Response *userIDResponseBody `json:"response"`
}

type userIDResponseBody struct {
	ID string `json:"id"`
}

type accountsResponse struct {
	Response *accountsResponseBody `json:"response"`
}

type accountsResponseBody struct {
	Accounts []Account `json:"accounts"`
}

type Account struct {
	ID       string  `json:"id"`
	IsActive bool    `json:"isActive"`
	Name     string  `json:"accountDisplayName"`
	Balance  float64 `json:"balance"`
}
