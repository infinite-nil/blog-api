package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-github/v57/github"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	app := echo.New()
	client := github.NewClient(nil)
	repoOptions := &github.RepositoryListByUserOptions{Type: "public"}

	app.GET("/healthy", func (ctx echo.Context) error {
		return ctx.String(http.StatusOK, "OK")
	})

	app.GET("/repositories", func(ctx echo.Context) error {
		repos, _, err := client.Repositories.ListByUser(context.Background(), "infinite-nil", repoOptions)

		if err != nil {
			// TODO: Improve error handling
			return echo.ErrInternalServerError
		}

		cleanRepos := make([]*github.Repository, len(repos))

		for i := 0; i < len(repos); i++ {
			cleanRepos[i] = &github.Repository{
				Name: repos[i].Name,
				HTMLURL: repos[i].HTMLURL,
				Language: repos[i].Language,
				Description: repos[i].Description,
				Topics: repos[i].Topics,
			}
		}

		return ctx.JSON(http.StatusOK, cleanRepos)
	})

	app.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus: true,
		LogURI:    true,
		BeforeNextFunc: func(c echo.Context) {
			c.Set("customValueFromContext", 42)
		},
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			value := c.Get("customValueFromContext")
			fmt.Printf("REQUEST: uri: %v, status: %v, custom-value: %v\n", v.URI, v.Status, value)
			return nil
		},
	}))

	app.Logger.Fatal(app.Start(":10000"))
}