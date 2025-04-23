package main

import (
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5"
	"github.com/pule1234/VideoForge/api"
	"github.com/pule1234/VideoForge/cache"
	"github.com/pule1234/VideoForge/cloud"
	"github.com/pule1234/VideoForge/config"
	db "github.com/pule1234/VideoForge/db/sqlc"
	"github.com/pule1234/VideoForge/gapi"
	"github.com/pule1234/VideoForge/global"
	"github.com/pule1234/VideoForge/internal/crawler"
	"github.com/pule1234/VideoForge/internal/processor"
	"github.com/pule1234/VideoForge/internal/publisher"
	"github.com/pule1234/VideoForge/mq"
	"github.com/pule1234/VideoForge/pb"
	"github.com/pule1234/VideoForge/worker"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

func main() {
	defer global.GlobalCancel()
	loadConfig, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}

	mq.InitRabbitMQ()
	cache.InitRedis()

	conn, err := pgx.Connect(global.GlobalCtx, loadConfig.DBSource)

	if err != nil {
		log.Fatal().Err(err).Msg("connect postgres err")
	}
	q := db.New(conn)
	//初始化QiNiu
	cloud.InitQiNiu(q)
	// 初始化 Publisher 工厂
	factory := publisher.NewPublisherFactory(q)
	dyCrawler, err := crawler.NewDyCrawler(loadConfig.DouYingQueueName, q)
	if err != nil {
		log.Error().Err(err).Msg("创建抖音爬虫失败")
	}

	log.Info().Msg("启动消息消费者...")
	go dyCrawler.Rabbit.ConsumeItem(processor.CreateCopyWriting, loadConfig.DouYingQueueName, dyCrawler.Postgres, global.GlobalCtx)

	redisOpt := asynq.RedisClientOpt{
		Addr: loadConfig.RedisSource,
	}
	taskDistributor := worker.NewTaskDistributor(redisOpt)
	taskprocessor := worker.NewTaskProcessor(redisOpt)
	go taskprocessor.Start()

	log.Info().Msg("启动Gin服务器...")
	go runGinServer(*loadConfig, q, factory, taskDistributor, taskprocessor)

	log.Info().Msg("启动gRPC服务器...")
	runGrpcServer(*loadConfig, q, factory, taskDistributor, taskprocessor)
}

func runGinServer(config config.Config, store db.Store, factory *publisher.PublisherFactory, taskDistributor *worker.TaskDistributor, taskprocessor *worker.TaskProcessor) {

	server, err := api.NewServer(config, store, factory, taskDistributor, taskprocessor)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create GIN server")
	}

	log.Info().Str("address", config.HTTPServerAddress).Msg("Gin服务器开始监听")
	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("GIN server failed to start")
	}
}

func runGrpcServer(
	config config.Config,
	store db.Store,
	factory *publisher.PublisherFactory,
	taskDistributor *worker.TaskDistributor,
	taskprocessor *worker.TaskProcessor,
) {
	server, err := gapi.NewServer(config, store, factory, taskDistributor, taskprocessor)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create gRPC server")
	}

	//集成自定义日志服务
	gprcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(gprcLogger)
	pb.RegisterVideosForgeServer(grpcServer, server)
	reflection.Register(grpcServer)

	listen, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener")
	}
	log.Info().Str("address", config.GRPCServerAddress).Msg("gRPC服务器开始监听")

	log.Info().Msg("gRPC服务器开始提供服务")
	err = grpcServer.Serve(listen)
	if err != nil {
		log.Error().Err(err).Msg("gRPC server failed to serve")
	}
	log.Info().Msg("gRPC服务器已关闭")
}
