package internal

import (
	"fmt"
	"gt-monitor/api"
	"gt-monitor/common"
	"gt-monitor/common/kafka"
	"gt-monitor/common/redis"
	"gt-monitor/common/zap"
	"gt-monitor/config"
	"gt-monitor/services"
	"gt-monitor/services/cache"
	"gt-monitor/utils"
	routing "github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"

	"log"
)

type ImpClkServer struct {
	cnf     *config.Config
	handler *services.ImpClkHandler
}

func NewImpClickServer(cnf *config.Config) *ImpClkServer {
	return &ImpClkServer{cnf: cnf}
}

func (s *ImpClkServer) Init() error {
	logger := zap.Get()
	logger.Info("start entity redis init")
	redisEntityPool, err := redis.NewPool(&s.cnf.Redis.Entity)
	if err != nil {
		logger.Error("entity redis init error", err)
		return err
	}

	logger.Info("start pubsub redis init")
	pubsubPool, err := redis.NewPool(&s.cnf.Redis.Pubsub)
	if err != nil {
		logger.Error("pubsub redis init error", err)
		return err
	}

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

	logger.Info("start entity cache service init")
	entityCache := cache.NewEntityCacheService(redisEntityPool)

	//
	logger.Info("start redis pub/sub init")
	flushLocalCacheFunc := func(data []byte) {
		logger.Info("received redis event: ", string(data))
		entityCache.FlushLocal()
	}

	pubSubService := redis.NewPubSub(pubsubPool)
	err = pubSubService.Subscribe("redis.data.track.link.channel", flushLocalCacheFunc)
	if err != nil {
		logger.Error("subscribe channel error ", err)
		return common.Redis_PubSub_Error
	}

	err = pubSubService.Subscribe("redis.data.creative.channel", flushLocalCacheFunc)
	if err != nil {
		logger.Error("subscribe channel error ", err)
		return common.Redis_PubSub_Error
	}

	logger.Info("start Kqueue init")
	kqueue := services.NewKqueue(producer)

	logger.Info("start handler init")
	handler := services.NewImpClkHandler(s.cnf.Kafka.Topic, entityCache, adverSyncer, kqueue)
	s.handler = handler

	logger.Info("done server init")
	return nil

}

func (s *ImpClkServer) Run() {
	logger := zap.Get()
	imp := api.NewImp(s.handler)
	click := api.NewClick(s.handler)

	router := routing.New()
	router.GET("/imp", imp.Handle)
	router.GET("/clk", click.Handle)
	router.GET("/clk/hello", test)

	address := fmt.Sprintf(": %d", s.cnf.Port)
	logger.Info("ImpClkServer bind to ", address)
	if err := fasthttp.ListenAndServe(address, api.CORS(router.Handler)); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}
func test(ctx *fasthttp.RequestCtx) {
	_, _ = ctx.WriteString("Hello Genius!")
}
