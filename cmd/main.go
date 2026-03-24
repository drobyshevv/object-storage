package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jackc/pgx/v5"

	"github.com/drobyshevv/object-storage/internal/db"
	"github.com/drobyshevv/object-storage/internal/handler"
	"github.com/drobyshevv/object-storage/internal/repository"
	"github.com/drobyshevv/object-storage/internal/service"
	"github.com/drobyshevv/object-storage/internal/storage"
	"github.com/drobyshevv/object-storage/pkg/config"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	var conn *pgx.Conn
	var err error
	for i := 0; i < 10; i++ {
		conn, err = db.New(cfg.DB.URL)
		if err == nil {
			break
		}
		log.Println("Waiting for DB, retrying in 2s...")
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatal("Cannot connect to DB:", err)
	}

	if err := db.RunMigrations(conn); err != nil {
		log.Fatal("Migration error:", err)
	}

	var client *s3.Client
	for i := 0; i < 10; i++ {
		awsCfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
			awsconfig.WithRegion(cfg.S3.Region),
			awsconfig.WithCredentialsProvider(
				credentials.NewStaticCredentialsProvider(
					cfg.S3.AccessKey,
					cfg.S3.SecretKey,
					"",
				),
			),
		)
		if err != nil {
			log.Println("AWS config error, retrying...", err)
			time.Sleep(2 * time.Second)
			continue
		}

		client = s3.NewFromConfig(awsCfg, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(cfg.S3.Endpoint)
			o.UsePathStyle = true
		})

		if err := storage.EnsureBucket(context.Background(), client, cfg.S3.Bucket); err != nil {
			log.Println("Bucket not ready, retrying...", err)
			time.Sleep(2 * time.Second)
			continue
		}

		break
	}
	if client == nil {
		log.Fatal("Cannot connect to S3")
	}

	repo := repository.NewFileRepository(conn)
	storageLayer := storage.New(client, cfg.S3.Bucket)
	serviceLayer := service.NewFileService(repo, storageLayer)
	handlerLayer := handler.NewFileHandler(serviceLayer)

	r := gin.Default()
	r.POST("/upload", handlerLayer.Upload)
	r.GET("/files", handlerLayer.List)
	r.GET("/files/:id", handlerLayer.Download)
	r.DELETE("/files/:id", handlerLayer.Delete)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Println("Server started on", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
