package response

type ConfigItem struct {
	Id      string `param:"id"`
	DataId  string `param:"dataId"`
	Group   string `param:"group"`
	Content string `param:"content"`
	Md5     string `param:"md5"`
	Tenant  string `param:"tenant"`
	Appname string `param:"appname"`
}
type ConfigPage struct {
	TotalCount     int          `param:"totalCount"`
	PageNumber     int          `param:"pageNumber"`
	PagesAvailable int          `param:"pagesAvailable"`
	PageItems      []ConfigItem `param:"pageItems"`
}

type ConfigListenContext struct {
	Group  string `json:"group"`
	Md5    string `json:"md5"`
	DataId string `json:"dataId"`
	Tenant string `json:"tenant"`
}

type ConfigContext struct {
	Group  string `json:"group"`
	DataId string `json:"dataId"`
	Tenant string `json:"tenant"`
}
