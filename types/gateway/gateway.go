package gateway

type GatewayData struct {
	Token string                 `json:"token"`
	User  string                 `json:"user"`
	More  map[string]interface{} `json:"more"`
}

type GatewayRecord struct {
	Data GatewayData `json:"data"`
	Url  string      `json:"url"`
}
