package min

import (
	"github.com/duanchi/min/bean"
	"github.com/duanchi/min/db/xorm"
	_interface "github.com/duanchi/min/interface"
	"github.com/duanchi/min/log"
	"github.com/gin-gonic/gin"
)

var HttpServer *gin.Engine
var Db *xorm.Engine
var Config interface{}
var Log *log.Logger
var DbMap map[string]*xorm.Engine
var Discovery _interface.DiscoveryInterface

func GetBean(name string) interface{} {
	return bean.Get(name)
}
