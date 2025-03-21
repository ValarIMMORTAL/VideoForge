package main

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/pule1234/VideoForge/api"
	"github.com/pule1234/VideoForge/cache"
	"github.com/pule1234/VideoForge/config"
	db "github.com/pule1234/VideoForge/db/sqlc"
	"github.com/pule1234/VideoForge/internal/crawler"
	"github.com/pule1234/VideoForge/internal/processor"
	"github.com/pule1234/VideoForge/mq"
	"log"
)

func main() {
	loadConfig, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config", err)
	}

	mq.InitRabbitMQ()
	cache.InitRedis()

	conn, err := pgx.Connect(context.Background(), loadConfig.DBSource)

	if err != nil {
		log.Fatal("connect postgres err ", err)
	}
	q := db.New(conn)

	dyCrawler, err := crawler.NewDyCrawler(loadConfig.DouYingQueueName, q)
	go dyCrawler.Rabbit.ConsumeItem(processor.CreateCopyWriting, loadConfig.DouYingQueueName, dyCrawler.Postgres)
	runGinServer(*loadConfig, q)
}

func runGinServer(config config.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create GIN server", err)
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("gRPC server failed to start", err)
	}
}

//func main() {
//	r := gin.Default()
//	r.GET("/ping", func(c *gin.Context) {
//		c.JSON(http.StatusOK, gin.H{
//			"message": "pong",
//		})
//	})
//	loadConfig, _ := config.LoadConfig(".")
//	r.Run(loadConfig.HTTPServerAddress) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
//}
