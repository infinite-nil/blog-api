package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	app := echo.New()

	app.GET("/", func(context echo.Context) error {
		return context.String(http.StatusOK, "OK")
	})

	app.Logger.Fatal(app.Start(":10000"))
}