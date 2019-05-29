package internal

import (
	"fmt"
	"gt-monitor/api"
	"gt-monitor/common/kafka"
	"gt-monitor/common/mongo"
	"gt-monitor/common/zap"
	"gt-monitor/config"
	"gt-monitor/services"
	"gt-monitor/services/dao"
	"gt-monitor/services/delayed"
	"gt-monitor/utils"
	routing "github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"log"
)

type ConvServer struct {
	cnf     *config.Config
	handler *services.ConvHandler
}

func NewConvServer(cnf *config.Config) *ConvServer {
	return &ConvServer{cnf: cnf}
}

func (s *ConvServer) Init() error {
	logger := zap.Get()

	logger.Info("init log  mongo db")
	logMgoPool, err := mongo.NewMgoPool(&s.cnf.Mongodb.Log)
	if err != nil {
		logger.Error("mongo pool init error", err)
		return err
	}

	logger.Info("init report  mongo db")
	reportMgoPool, err := mongo.NewMgoPool(&s.cnf.Mongodb.Report)
	if err != nil {
		logger.Error("mongo pool init error", err)
		return err
	}

	logger.Info("start mongo dao init")
	mongoDao := dao.NewMongoDao(logMgoPool, reportMgoPool)

	logger.Info("start kafka init")
	producer, err := kafka.New(&s.cnf.Kafka.Server)
	if err != nil {
		logger.Error("kafka init error", err)
		return err
	}

	logger.Info("start http client init")
	httpCli := utils.NewHttpCli(&s.cnf.Httpcli)

	logger.Info("start adver syncer init")
	adverSyncer := services.NewAdverSyncer(httpCli)

	logger.Info("start Kqueue init")
	kqueue := services.NewKqueue(producer)

	logger.Info("start task queue init dir=", s.cnf.TaskQueueDataDir)
	delayedTaskQueue, err := delayed.NewDelayedQueue(s.cnf.TaskQueueDataDir)
	if err != nil {
		logger.Error("delayed task queue init error", err)
		return err
	}

	logger.Info("start handler init")
	handler := services.NewConvHandler(s.cnf.Kafka.Topic, mongoDao, adverSyncer, kqueue, delayedTaskQueue)
	s.handler = handler

	convWorker := services.NewConvWorker(handler, delayedTaskQueue, mongoDao)
	convWorker.Run()

	logger.Info("done server init")
	return nil

}

func (s *ConvServer) Run() {
	logger := zap.Get()
	conv := api.NewConversion(s.handler)

	router := routing.New()
	router.GET("/e", conv.Handle)
	router.POST("/e", conv.Handle)
	router.GET("/e/hello/:name", Hello)

	address := fmt.Sprintf(": %d", s.cnf.Port)
	logger.Info("ConvServer bind to ", address)
	if err := fasthttp.ListenAndServe(address, api.CORS(router.Handler)); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}

func Hello(ctx *fasthttp.RequestCtx) {
	_, _ = fmt.Fprintf(ctx, "hello, %s!\n", ctx.UserValue("name"))
}
