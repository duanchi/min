package core_parsers

import (
	"github.com/duanchi/min/v2/event"
	_interface "github.com/duanchi/min/v2/interface"
	"github.com/duanchi/min/v2/rpc"
	"github.com/duanchi/min/v2/scheduled"
	"github.com/duanchi/min/v2/server/middleware"
	"github.com/duanchi/min/v2/server/route"
	"github.com/duanchi/min/v2/server/validate"
	"github.com/duanchi/min/v2/service"
)

var CoreBeanParsers = []_interface.BeanParserInterface{
	&service.ServiceBeanParser{},
	&route.RouteBeanParser{},
	&route.RestfulBeanParser{},
	&middleware.MiddlewareBeanParser{},
	&scheduled.ScheduledBeanParser{},
	&event.EventBeanParser{},
	&rpc.RpcBeanParser{},
	&validate.ValidatorBeanParser{},
}
