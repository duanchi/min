package min

import (
	"fmt"
	"github.com/duanchi/min/bean"
	"github.com/duanchi/min/cache"
	"github.com/duanchi/min/config"
	"github.com/duanchi/min/db"
	"github.com/duanchi/min/log"
	"github.com/duanchi/min/microservice/discovery"
	"github.com/duanchi/min/scheduled"
	"github.com/duanchi/min/server"
	config2 "github.com/duanchi/min/types/config"
	"github.com/joho/godotenv"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	initEnv()
}

func Bootstrap(configuration interface{}) {
	config.Init(configuration)
	Config = configuration

	errs := make(chan error, 3)

	bean.InitBeans(
		config.Get("Beans"),
		config.Get("BeanParsers"),
	)

	log.Init(config.Get("Log").(config2.Log))
	Log = &log.Log
	if !checkConfigEnabled("Log.Enabled") {
		Log.Enabled(false)
	}

	if checkConfigEnabled("Scheduled.Enabled") {
		Log.Info("Task Enabled!")
		scheduled.Init()
	}

	if checkConfigEnabled("Db.Enabled") {
		db.Init()
		Db = db.Connection
	}

	if checkConfigEnabled("Cache.Enabled") {
		cache.Init()
	}

	if checkConfigEnabled("Discovery.Enabled") {
		discovery.Init()
		Discovery = discovery.Discovery
	}

	/*if checkConfigEnabled("Feign.Enabled") {
		feign.Init(config.Get("Feign").(config2.Feign))
	}*/

	if checkConfigEnabled("Scheduled.Enabled") {
		go scheduled.RunOnInit()
	}

	go server.Init(errs)
	HttpServer = server.HttpServer

	if checkConfigEnabled("Scheduled.Enabled") {
		go scheduled.RunOnStart()
	}

	go func() {
		c := make(chan os.Signal, 2)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	err := <-errs

	if checkConfigEnabled("Scheduled.Enabled") {
		go scheduled.RunOnExit()
	}

	log.Log.Error("%s", err)
}

func SetConfigFile(configFile string) {
	config.SetConfigFile(configFile)
}

func checkConfigEnabled(configStack string) bool {
	return config.Get(configStack).(bool)
}

func initEnv() {
	env := os.Getenv("ENV")
	envFile := ".env."
	switch env {
	case "production":
		envFile += "prod"
	case "test":
		envFile += "test"
	case "development":
		fallthrough
	default:
		envFile += "dev"
	}
	fmt.Println("load env file " + envFile + ".local")
	fmt.Println("load env file " + ".env.local")
	fmt.Println("load env file " + ".env")
	godotenv.Overload(".env", ".env.local", envFile+".local")
}
