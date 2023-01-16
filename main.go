package main

import (
	"cpm-rad-backend/domain/config"
	"cpm-rad-backend/domain/connection"
	"cpm-rad-backend/domain/logger"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
)

func main() {
	config.InitConfig()

	// Zaplog
	zaplog, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialze zap logger : %v", err)
	}
	defer zaplog.Sync()

	// Database
	// cpmDB, err := gorm.Open(sqlserver.Open(config.DBCpm), &gorm.Config{})
	// if err != nil {
	// 	log.Fatalf("can't connect DB : %v", err)
	// }

	// db := &connection.DBConnection{
	// 	CPM: cpmDB,
	// }

	db := &connection.DBConnection{}

	// Migrate the schema
	// db.AutoMigrate(&Product{})

	e := getRoute(zaplog)

	// Routes
	InitAPIV1(e.Group("/api/v1"), db)

	e.Logger.Fatal(e.Start(":8000"))

}

func getRoute(zaplog *zap.Logger) *echo.Echo {
	// Echo instance
	e := echo.New()

	e.Logger.SetLevel(log.INFO)
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{}))
	e.Use(logger.Middleware(zaplog))
	return e
}

func InitAPIV1(api *echo.Group, db *connection.DBConnection) {

	fmt.Print(db)

	api.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!!")
	})
}
