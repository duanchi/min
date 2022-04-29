package request

type Listener func(namespace, group, dataId, data string)

type Config struct {
	DataId           string `param:"dataId"`  //required
	Group            string `param:"group"`   //required
	Content          string `param:"content"` //required
	Tag              string `param:"tag"`
	AppName          string `param:"appName"`
	BetaIps          string `param:"betaIps"`
	CasMd5           string `param:"casMd5"`
	Type             string `param:"type"`
	EncryptedDataKey string `param:"encryptedDataKey"`
	OnChange         func(namespace, group, dataId, data string)
}

type SearchConfig struct {
	Search   string `param:"search"`
	DataId   string `param:"dataId"`
	Group    string `param:"group"`
	Tag      string `param:"tag"`
	AppName  string `param:"appName"`
	PageNo   int    `param:"pageNo"`
	PageSize int    `param:"pageSize"`
}
