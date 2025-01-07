package core_parsers

import (
	"github.com/duanchi/min/event"
	_interface "github.com/duanchi/min/interface"
	"github.com/duanchi/min/rpc"
	"github.com/duanchi/min/scheduled"
	"github.com/duanchi/min/server/middleware"
	"github.com/duanchi/min/server/route"
	"github.com/duanchi/min/server/validate"
	"github.com/duanchi/min/service"
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
