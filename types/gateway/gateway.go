package gateway

type Data struct {
	Data struct{
		Token string `json:"token"`
		User string `json:"user"`
		More map[string]interface{} `json:"more"`
	} `json:"data"`
	Url string `json:"url"`
}
