package holder

import (
	"github.com/duanchi/min/microservice/discovery/nacos/request"
	"github.com/duanchi/min/requests/http"
	"strconv"
)

var leader = 0

type HttpHolder struct {
	requestHolder *http.Request
}

func NewHttpHolder(client request.Client) (httpHolder HttpHolder) {
	if leader > -1 && len(client.ServerConfigs) > leader {
		service := client.ServerConfigs[leader]

		holder := http.New()
		holder.BaseUrl(service.Scheme + "://" + service.IpAddr + ":" + strconv.FormatUint(service.Port, 10) + service.ContextPath)
		httpHolder = HttpHolder{requestHolder: &holder}
	}

	return
}

func (this *HttpHolder) GET(url string, parameters interface{}, response interface{}) (err error) {
	responseData, err := this.requestHolder.Method("GET").Query(parameters).Response()

	if err == nil {
		err = responseData.BindJSON(response)
	}

	return
}

func (this *HttpHolder) POST(url string, parameters interface{}, response interface{}) (ok bool, err error) {
	_, err = this.requestHolder.Method("POST").Form(parameters).Response()

	if err == nil {
		ok = true
	}

	return
}

func (this *HttpHolder) PUT(url string, parameters interface{}, response interface{}) (ok bool, err error) {
	_, err = this.requestHolder.Method("PUT").Form(parameters).Response()

	if err == nil {
		ok = true
	}

	return
}

func (this *HttpHolder) DELETE(url string, parameters interface{}) (ok bool, err error) {
	_, err = this.requestHolder.Method("DELETE").Query(parameters).Response()

	if err == nil {
		ok = true
	}

	return
}

func getLeader() int {
	return leader
}
