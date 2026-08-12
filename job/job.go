package job

import (
	"context"
	"encoding/json"
	"net/url"
	"reflect"

	config2 "github.com/duanchi/min/v2/config"
	_interface "github.com/duanchi/min/v2/interface"
	"github.com/duanchi/min/v2/job/xxl"
	"github.com/duanchi/min/v2/log"
	"github.com/duanchi/min/v2/types/config"
	"github.com/duanchi/min/v2/util"
)

var JobList = map[string]reflect.Value{}

var exec xxl.Executor

var cancel context.CancelFunc

func Init() {

	if len(JobList) == 0 {
		return
	}

	jobConfig := config2.Get("Job").(config.Job)
	applicationConfig := config2.Get("Application").(config.Application)
	executorIp := jobConfig.ExecutorIp
	executorPort := jobConfig.ExecutorPort

	if executorIp == "" {
		executorIp = util.GetIp()
	}

	if executorPort == "" {
		executorPort = "9999"
	}

	u, _ := url.Parse(jobConfig.Server)
	u.Path = "/xxl-job-admin"

	opts := []xxl.Option{
		xxl.ServerAddr(u.String()),
		xxl.ExecutorIp(executorIp),              //可自动获取
		xxl.ExecutorPort(executorPort),          //默认9999（非必填）
		xxl.RegistryKey(applicationConfig.Name), //执行器名称
		// xxl.SetLogger(&logger{}),                     //自定义日志
	}
	if jobConfig.AccessToken != "" {
		opts = append(opts, xxl.AccessToken(jobConfig.AccessToken))
	}
	exec = xxl.NewExecutor(opts...)
	exec.Init()
	if jobConfig.EnableRegister {
		exec.Register()
	}
	//设置日志查看handler
	// exec.LogHandler(customLogHandle)
	//注册任务handler

	for name, job := range JobList {
		exec.RegTask(name, func(cxt context.Context, param *xxl.RunReq) string {
			executorParams := map[string]any{}
			json.Unmarshal([]byte(param.ExecutorParams), &executorParams)
			err := job.Interface().(_interface.JobInterface).Execute(executorParams, param.ExecutorParams)

			if err != nil {
				return err.Error()
			}
			return "OK"
		})
	}
	ctx, c := context.WithCancel(context.Background())
	defer c()
	cancel = c
	err := exec.Run(ctx)
	if err != nil {
		log.Log.Error("[min-framework]: Job run error, \"" + err.Error() + "\"")
	}
}

func Stop() {
	if cancel != nil {
		cancel()
	}
}
