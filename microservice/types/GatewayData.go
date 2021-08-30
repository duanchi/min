package types

type GatewayData struct {
	Data struct{
		Token string `json:"token"`
		User string `json:"User"`
		More string `json:"more"`
	} `json:"data"`
	Url string `json:"url"`
}
