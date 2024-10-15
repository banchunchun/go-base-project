package main

import (
	"com.banxiaoxiao.server/config"
	"com.banxiaoxiao.server/etcd"
	"com.banxiaoxiao.server/logger"
	"com.banxiaoxiao.server/middleware"
	"com.banxiaoxiao.server/migration"
	"com.banxiaoxiao.server/model"
	"com.banxiaoxiao.server/repo"
	"com.banxiaoxiao.server/repository"
	"com.banxiaoxiao.server/router"
	"com.banxiaoxiao.server/util"
	"context"
	"flag"
	"github.com/labstack/echo/v4"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var cfgPath string
	var workPath string
	flag.StringVar(&cfgPath, "c", "", "config path")
	flag.StringVar(&workPath, "w", "", "work path")
	flag.Parse()

	util.ChangeWorkPath(workPath)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	//appCtx, cancelApp := context.WithCancel(context.Background())
	//defer cancelApp()

	conf := config.Load(cfgPath)

	model.InitNode()
	l := logger.NewLogger(cfgPath)
	logger.Log().Infof("start main application")

	e := echo.New()

	rep := repository.NewProxyRepository(conf, l)
	defer rep.Close()
	repo.SetRepository(rep)
	migration.CreateDatabase()

	etcdClient, _ := etcd.NewEtcd()
	if etcdClient != nil {
		etcdClient.Leader()
	}

	router.Init(e)
	middleware.InitLoggerMiddleware(e)

	go func() {
		if err := e.Start(conf.Web.Host + ":" + conf.Web.Port); err != nil {
			logger.Log().Info(err.Error())
			stop()
		}
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
