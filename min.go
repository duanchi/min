package min

import (
	"github.com/duanchi/min/v2/context"
	"github.com/duanchi/min/v2/db/xorm"
	"github.com/duanchi/min/v2/log"
	"github.com/duanchi/min/v2/server/httpserver"
)

var HttpServer *httpserver.Httpserver
var Db *xorm.Engine
var Config interface{}
var Log *log.Logger
var DbMap map[string]*xorm.Engine
var ApplicationContext = context.GetApplicationContext()
