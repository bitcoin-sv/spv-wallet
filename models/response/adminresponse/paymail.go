package adminresponse

type Paymail struct {
	ID uint `json:"id"`

	Alias   string `json:"alias"`
	Domain  string `json:"domain"`
	Paymail string `json:"paymail"`

	PublicName string `json:"publicName"`
	Avatar     string `json:"avatar"`
}
