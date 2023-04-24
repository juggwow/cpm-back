package main

import (
	"context"
	"cpm-rad-backend/domain/auth"
	"cpm-rad-backend/domain/boq"
	"cpm-rad-backend/domain/config"
	"cpm-rad-backend/domain/connection"
	"cpm-rad-backend/domain/contract"
	"cpm-rad-backend/domain/employee"
	"cpm-rad-backend/domain/form"
	"cpm-rad-backend/domain/health_check"
	"cpm-rad-backend/domain/logger"
	"cpm-rad-backend/domain/minio"
	"cpm-rad-backend/domain/raddoc"
	"cpm-rad-backend/domain/report"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
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

	minioClient := initMinio()

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

func initMinio() minio.Client {
	conf := minio.Configuration{
		Endpoint:   config.StorageEndpoint,
		AccessKey:  config.StorageAccessKey,
		SecretKey:  config.StorageSecretKey,
		UseSSL:     config.StorageSSL,
		BucketName: config.StorageBucketName,
	}
	if err := minio.NewConnection(conf); err != nil {
		log.Fatalf("can't connect MINIO client: %v", err)
		panic(err)
	}

	minioClient := minio.GetClient()
	return minioClient
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

func getAuthMiddleware() echo.MiddlewareFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:      &auth.JwtEmployeeClaims{},
		SigningKey:  []byte(config.AuthJWTSecret),
		TokenLookup: "header:Authorization,cookie:" + config.AuthJWTKey,
	})
}

func initPublicAPI(e *echo.Echo, db *connection.DBConnection, minioClient minio.Client) {

	e.GET("/healths", health_check.HealthCheck)
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	if authenticator, err := auth.NewAuthenticator(); err == nil {
		e.GET("/auth", authenticator.AuthenHandler())
		e.GET("/auth/callback", authenticator.AuthenCallbackHandler(employee.GetAndCreateIfNotExist(db), auth.CreateLog(db)))
		e.GET("/auth/refreshToken", authenticator.GetRefreshTokenHandler(auth.GetIDToken(db)), getAuthMiddleware())
		e.GET("/auth/logout/:token", authenticator.LogoutHandler(auth.GetIDToken(db)))

	} else {
		log.Fatalf("Fatal initiate authenticator: %v\n", err)
		panic(err)
	}
}

func initAPIV1(api *echo.Group, db *connection.DBConnection, minioClient minio.Client) {
	if config.AuthJWTEnabled {
		api.Use(getAuthMiddleware())
	}
	api.GET("/employees/me", auth.GetCurrentHandler)
	api.GET("/employees/:employeeId", employee.GetByIDHandler(employee.GetByID(db)))

	api.GET("/contract/:id", contract.GetByIDHandler(contract.GetByID(db)))
	api.GET("/contract/:id/boq", boq.GetHandler(boq.Get(db)))
	api.GET("/contract/:id/card", contract.GetNumberOfItemHandler(contract.GetNumberOfItem(db)))

	api.GET("/boq/:id", boq.GetItemByIDHandler(boq.GetItemByID(db)))
	api.GET("/country", form.GetCountryHandler(form.GetCountry(db)))
	api.GET("/doctype", form.GetDocTypeHandler(form.GetDocType(db)))
	api.POST("/form", report.CreateHandler(report.Create(db, minioClient)))
	api.GET("/form/:id", form.GetHandler(form.Get(db)))
	api.PUT("/form/:id", report.UpdateHandler(report.Update(db, minioClient)))
	api.DELETE("/form/:id", form.DeleteHandler(form.Delete(db)))

	// api.POST("/upload/:fieldName/:itemid", form.FileUploadHandler(form.FileUpload(db, minioClient)))
	// api.DELETE("/delete/:itemid", form.FileDeleteHandler(form.FileDelete(db, minioClient)))
	api.GET("/download/:fileid", form.FileDownloadHandler(form.FileDownload(db, minioClient)))

	api.GET("/listofdoc/:itemid", raddoc.GetByItemHandler(raddoc.GetByItem(db)))
	api.GET("/listofdoc/progress/contract/:id", report.GetProgressReportHandler(report.GetProgressReport(db)))
	api.GET("/listofdoc/check/contract/:id", report.GetCheckReportHandler(report.GetCheckReport(db)))

}
