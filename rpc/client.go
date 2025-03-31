package rpc

import (
	"bytes"
	"encoding/json"
	_interface "github.com/duanchi/min/v2/interface"
	"github.com/duanchi/min/v2/types"
	"github.com/duanchi/min/v2/util"
	"io/ioutil"
	"net/http"
	"reflect"
	"runtime"
	"strings"
)

func Call(in IN, out *OUT, caller _interface.RpcInterface) (err error) {
	pc, _, _, _ := runtime.Caller(1)
	methodStack := strings.Split(runtime.FuncForPC(pc).Name(), ".")

	err = rpcRequest(
		caller.GetApplicationName(),
		caller.GetPackageName(),
		reflect.ValueOf(caller).Elem().Type().String(),
		methodStack[len(methodStack)-1],
		&in,
		out,
	)

	return
}

func rpcRequest(serviceName string, packageName string, className string, method string, in *IN, out *OUT) (err error) {
	client := &http.Client{}

	requestBody, _ := json.Marshal(in)

	request, err := http.NewRequest(http.MethodPost, "http://"+serviceName+"/"+packageName+"/"+className+"::"+method, bytes.NewReader(requestBody))
	if err != nil {
		// handle error
		err = types.RuntimeError{
			Message:   "RPC request error",
			ErrorCode: 500,
		}
		return
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Request-Id", util.GenerateUUID().String())

	response, err := client.Do(request)

	defer response.Body.Close()

	if err != nil {
		err = types.RuntimeError{
			Message:   "RPC response error, " + err.Error(),
			ErrorCode: 500,
		}
		return
	}

	responseBody, err := ioutil.ReadAll(response.Body)

	err = json.Unmarshal(responseBody, out)

	return
}
