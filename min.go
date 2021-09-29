package min

import (
	"github.com/duanchi/min/bean"
	"github.com/duanchi/min/log"
	"github.com/gin-gonic/gin"
	"github.com/xormplus/xorm"
)

var HttpServer *gin.Engine
var Db *xorm.Engine
var Config interface{}
var Log *log.Logger
var DbMap map[string]*xorm.Engine

func GetBean(name string) interface{} {
	return bean.Get(name)
}


