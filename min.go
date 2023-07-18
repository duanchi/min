package min

import (
	"github.com/duanchi/min/context"
	"github.com/duanchi/min/db/xorm"
	"github.com/duanchi/min/log"
	"github.com/duanchi/min/server/httpserver"
)

var HttpServer *httpserver.Httpserver
var Db *xorm.Engine
var Config interface{}
var Log *log.Logger
var DbMap map[string]*xorm.Engine
var ApplicationContext = context.GetApplicationContext()
