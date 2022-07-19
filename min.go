package min

import (
	"github.com/duanchi/min/context"
	"github.com/duanchi/min/db/xorm"
	"github.com/duanchi/min/log"
	"github.com/gin-gonic/gin"
)

var HttpServer *gin.Engine
var Db *xorm.Engine
var Config interface{}
var Log *log.Logger
var DbMap map[string]*xorm.Engine
var ApplicationContext = context.NewApplicationContext()
