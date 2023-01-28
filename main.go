package main

import (
	"context"
	"cpm-rad-backend/domain/boq"
	"cpm-rad-backend/domain/config"
	"cpm-rad-backend/domain/connection"
	"cpm-rad-backend/domain/contract"
	"cpm-rad-backend/domain/form"
	"cpm-rad-backend/domain/health_check"
	"cpm-rad-backend/domain/logger"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/zap"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func main() {
	config.InitConfig()

	zaplog, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialze zap logger : %v", err)
	}
	defer zaplog.Sync()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	e := getRoute(zaplog)

	cpmDB, err := gorm.Open(sqlserver.Open(config.DBCpm), &gorm.Config{})
	if err != nil {
		log.Fatalf("can't connect DB : %v", err)
		panic(err)
	}

	db := &connection.DBConnection{
		CPM: cpmDB,
	}

	// db := &connection.DBConnection{}

	// err = cpmDB.AutoMigrate()

	// if err != nil {
	// 	return
	// }

	minioClient, err := minio.New(config.StorageEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.StorageAccessKey, config.StorageSecretKey, ""),
		Secure: config.StorageSSL,
	})
	if err != nil {
		log.Fatalf("can't connect MINIO client: %v", err)
		panic(err)
	}

	initPublicAPI(e, db, minioClient)

	initAPIV1(e.Group("/api/v1"), db, minioClient)

	go func() {
		if err := e.Start(":" + config.AppPort); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

}

func getRoute(zaplog *zap.Logger) *echo.Echo {
	e := echo.New()
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: 10 * time.Second,
	}))

	e.Logger.SetLevel(log.INFO)

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: func(c echo.Context) bool {
			return path.Base(c.Request().URL.Path) == "healths"
		},
	}))
	e.Use(middleware.Recover())
	e.Use(logger.Middleware(zaplog))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))

	return e
}

func initPublicAPI(e *echo.Echo, db *connection.DBConnection, minioClient *minio.Client) {
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!!")
	})
	e.GET("/db", func(c echo.Context) error {
		return c.String(http.StatusOK, fmt.Sprintf("Hello, DB!! >>> %v", config.DBCpm))
	})

	e.GET("/healths", health_check.HealthCheck)
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// e.GET("/uploads/:filename", upload.DownloadHandler(minioClient))

	// if authenticator, err := auth.NewAuthenticator(); err == nil {
	// 	e.GET("/auth", authenticator.AuthenHandler())
	// 	e.GET("/auth/callback", authenticator.AuthenCallbackHandler(employee.GetAndCreateIfNotExist(db)))
	// } else {
	// 	log.Fatalf("Fatal initiate authenticator: %v\n", err)
	// 	panic(err)
	// }
}

func initAPIV1(api *echo.Group, db *connection.DBConnection, minioClient *minio.Client) {

	//fmt.Print(db)
	api.GET("/contract/:id", contract.GetByIDHandler(contract.GetByID(db)))
	api.GET("/contract/:id/boq", boq.GetHandler(boq.Get(db)))

	api.GET("/boq/:id", boq.GetItemByIDHandler(boq.GetItemByID(db)))
	api.GET("/country", form.GetCountryHandler(form.GetCountry(db)))
	api.POST("/form", form.CreateHandler(form.Create(db)))

}
