package external

type boolResponse struct {
	Response *bool `json:"response"`
}

type stringResponse struct {
	Response *string `json:"response"`
}

type userIDResponse struct {
	Response *UserIDResponseBody `json:"response"`
}

type UserIDResponseBody struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
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
	Tender   string  `json:"accountTender"`
	Balance  float64 `json:"balance"`
}

type userPhotoResponse struct {
	Response *userPhotoResponseBody `json:"response"`
}

type userPhotoResponseBody struct {
	MimeType string `json:"mimeType"`
	Data     string `json:"data"`
}

type barcodeSeedResponse struct {
	Response *barcodeSeedResponseBody `json:"response"`
}

type barcodeSeedResponseBody struct {
	BarcodeSeed string `json:"barcodeSeed"`
}

type cashlessKeyResponse struct {
	Response *cashlessKeyResponseBody `json:"response"`
}

type cashlessKeyResponseBody struct {
	Value string `json:"value"`
}

type retrieveSettingResponse struct {
	Response retrieveSettingResponseBody `json:"response"`
}

type retrieveSettingResponseBody struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}
