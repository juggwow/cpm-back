package main

import (
	"cpm-rad-backend/domain/config"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func main() {
	config.InitConfig()

	fmt.Printf(config.DBCpm)
	db, err := gorm.Open(sqlserver.Open(config.DBCpm), &gorm.Config{})
	// Migrate the schema
	if err != nil {
		panic("failed to connect database")
	}

	fmt.Print(db)
	// db.AutoMigrate(&Product{})

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hello)

	e.Logger.Fatal(e.Start(":1323"))

}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
