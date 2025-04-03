package holder

import (
	"github.com/duanchi/min/v2/microservice/discovery/nacos/request"
	"github.com/duanchi/min/v2/requests/http"
	"strconv"
)

var leader = 0

type HttpHolder struct {
	requestHolder *http.Request
}

func NewHttpHolder(client request.Client) (httpHolder *HttpHolder) {
	if leader > -1 && len(client.ServerConfigs) > leader {
		service := client.ServerConfigs[leader]
		holder := http.New()
		holder.BaseUrl(service.Scheme + "://" + service.IpAddr + ":" + strconv.FormatUint(service.Port, 10) + service.ContextPath)
		httpHolder = &HttpHolder{requestHolder: &holder}
	}

	return
}

func (this *HttpHolder) GET(url string, parameters interface{}, response interface{}) (err error) {
	responseData, err := this.
		requestHolder.Url(url).Method("GET").Query(parameters).Response()

	if err == nil {
		err = responseData.BindJSON(response)
	}

	return
}

func (this *HttpHolder) POST(url string, parameters interface{}) (ok bool, err error) {
	_, err = this.requestHolder.Url(url).Method("POST").Form(parameters).Response()

	if err == nil {
		ok = true
	}

	return
}

func (this *HttpHolder) PUT(url string, parameters interface{}) (ok bool, err error) {
	_, err = this.requestHolder.Url(url).Method("PUT").Form(parameters).Response()

	if err == nil {
		ok = true
	}

	return
}

func (this *HttpHolder) Holder() *http.Request {
	return this.requestHolder
}

func (this *HttpHolder) DELETE(url string, parameters interface{}) (ok bool, err error) {
	_, err = this.requestHolder.Url(url).Method("DELETE").Query(parameters).Response()

	if err == nil {
		ok = true
	}

	return
}

func getLeader() int {
	return leader
}
