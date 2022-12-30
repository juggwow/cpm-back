package main

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func main() {
	dsn := "sqlserver://SA:" + url.QueryEscape("S#123456") + "@localhost:1433?database=CPM"
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
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
