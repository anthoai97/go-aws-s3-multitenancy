package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	api "github.com/anthoai97/go-aws-s3-multitenancy/api"
	"github.com/anthoai97/go-aws-s3-multitenancy/business"
	"github.com/anthoai97/go-aws-s3-multitenancy/core"
	md "github.com/anthoai97/go-aws-s3-multitenancy/middleware"
	"github.com/anthoai97/go-aws-s3-multitenancy/repository/storage_s3"
	"github.com/anthoai97/go-aws-s3-multitenancy/repository/token_vendor_machine"
	logger "github.com/ethereum/go-ethereum/log"
	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
)

type Env string

var (
	Log            = logger.New("logscope", "master")
	verbose        = flag.Bool("v", true, "more verbose logs")
	Develop    Env = "develop"
	Production Env = "production"
)

func init() {
	env := core.GetEnvVar("ENV", "develop")
	if env == string(Develop) {
		err := godotenv.Load()
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	flag.Parse()

	hs := logger.StreamHandler(os.Stderr, logger.TerminalFormat(true))
	loglevel := logger.LvlInfo
	if *verbose {
		loglevel = logger.LvlTrace
	}
	hf := logger.LvlFilterHandler(loglevel, hs)
	logger.Root().SetHandler(hf)
}

func main() {
	Log.Info("Start Storage Service....")
	router := gin.New()
	router.Use(gin.Recovery(), gin.Logger())

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "pong"})
	})

	setupRoutes(router)

	addr := fmt.Sprintf("0.0.0.0:%v", 8080)
	if err := router.Run(addr); err != nil {
		log.Fatal("Server run failed ", err)
	}
}

func setupRoutes(router *gin.Engine) {
	fmt.Println("Hear we go again")
	roleArn := "arn:aws:iam::360307311296:role/ri.developer-assume-role" // will be taken from config file
	externalID := "someexternalid"                                       // Must match role trust entity
	// tenant := "dataspire"
	bucket := core.GetEnvVar("BUCKET", "dsr-customer-storage-dev")

	// Middleware
	// Token vendor machine
	credDurationSeconds := core.GetEnvVar[int]("CRED_DURATION", 900)
	logger.Info("STS Credential Duration Seconds", "seconds", credDurationSeconds)
	tvm := token_vendor_machine.NewTokenVendorMachine(roleArn, externalID, bucket, int32(credDurationSeconds), logger.New("TokenVendorMachine", externalID))
	storage_s3 := storage_s3.NewStorageS3(bucket, logger.New("log-scope", "storage_s3"))

	// Business, API
	bussiness := business.NewBusiness(tvm, storage_s3, logger.New("log-scope", "business"))
	api := api.NewAPI(bussiness)

	router.POST("/credential", api.GenerateSTSCredentialHdl())

	routerS3 := router.Group("/s3")
	var middleware md.CustomMiddleware
	routerS3.Use(middleware.CheckSTSCrendential(logger.New("log-scope", "middleware")))
	{
		routerS3.GET("/tree", api.ListS3StorageTreeHdl())
		routerS3.POST("/files", api.UploadS3ObjectsByGenerateUrlHdl())
		// routerS3.POST("/folder", api.PostS3ObjectFolderHdl())
		// routerS3.POST("/files", api.PostS3FilesByGenerateUrlHdl())
		// routerS3.POST("/download/files", api.PostDownloadS3FilesHdl())
		// routerS3.POST("/download/folder", api.PostDownloadS3FolderHdl())
		// routerS3.DELETE("/objects", api.DeleteS3ObjectsHdl())
	}
}
