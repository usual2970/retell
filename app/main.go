package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/usual2970/retell/internal/rest/middleware"
	"github.com/usual2970/retell/pkg/config"
	"github.com/usual2970/retell/pkg/db"
	"github.com/usual2970/retell/pkg/logger"
	"github.com/usual2970/retell/pkg/mq/kafka"
	"github.com/usual2970/retell/pkg/redis"
	"github.com/usual2970/retell/pkg/redlock"
	"github.com/usual2970/retell/pkg/schedule"
	"github.com/usual2970/retell/pkg/sse"
	"github.com/usual2970/retell/pkg/validator"
	"github.com/usual2970/retell/routes"

	"github.com/labstack/echo-contrib/echoprometheus"

	echopprof "github.com/sevenNt/echo-pprof"
)

const (
	defaultTimeout = 30
	defaultAddress = ":9090"
)

func main() {
	// 初始化配置
	if err := config.Init(); err != nil {
		log.Fatal(err)
	}

	// 初始化日志
	if err := logger.Init(); err != nil {
		log.Fatal(err)
	}

	// 初始化数据库
	if err := db.Init(); err != nil {
		log.Fatal(err)
	}

	// 初始化redis
	if err := redis.Init(); err != nil {
		log.Fatal(err)
	}

	// 初始化锁
	redlock.InitLock()

	logger.Info("Starting the application...")

	// 获取echo
	e := getEcho()

	e.Use(echoprometheus.NewMiddleware("paas"))
	e.GET("/metrics", echoprometheus.NewHandler())

	// 注册sse
	sse.InitPushlet(e, "/pushlet/events")

	// 注册校验器
	validator.Register(e)

	// 注册路由
	routes.Register(e)

	// 初始化KAFKA
	kafka.SetupDefaultClient()

	// 注册定时任务
	schedule.StartSchedule()

	// Start Server
	address := os.Getenv("SERVER_ADDRESS")
	if address == "" {
		address = defaultAddress
	}
	if err := e.Start(address); err != nil {
		schedule.StopSchedule()
		log.Fatal(err)
	}

}

func getEcho() *echo.Echo {
	e := echo.New()

	// pprof
	echopprof.Wrap(e)

	e.Use(middleware.CORS)
	timeoutStr := os.Getenv("CONTEXT_TIMEOUT")
	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil {
		log.Println("failed to parse timeout, using default timeout")
		timeout = defaultTimeout
	}
	timeoutContext := time.Duration(timeout) * time.Second
	e.Use(middleware.SetRequestContextWithTimeout(timeoutContext))

	e.Use(middleware.Recover)

	return e
}
