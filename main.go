package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/google/go-github/v57/github"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	sentry.Init(sentry.ClientOptions{
		Dsn:              os.Getenv("SENTRY_DSN"),
		TracesSampleRate: 1.0,
	})

	app := echo.New()
	client := github.NewClient(nil)
	repoOptions := &github.RepositoryListByUserOptions{Type: "public"}

	app.GET("/healthy", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "OK")
	})

	app.GET("/repositories", func(ctx echo.Context) error {
		repos, _, err := client.Repositories.ListByUser(context.Background(), "infinite-nil", repoOptions)

		if err != nil {
			// TODO: Improve error handling
			log.Fatal("GIHUB ERROR", err)
			return echo.ErrInternalServerError
		}

		cleanRepos := make([]*github.Repository, len(repos))

		for i := 0; i < len(repos); i++ {
			cleanRepos[i] = &github.Repository{
				Name:        repos[i].Name,
				HTMLURL:     repos[i].HTMLURL,
				Language:    repos[i].Language,
				Description: repos[i].Description,
				Topics:      repos[i].Topics,
			}
		}

		return ctx.JSON(http.StatusOK, cleanRepos)
	})

	app.Use(middleware.Recover())
	app.Use(middleware.Logger())
	app.Use(sentryecho.New(sentryecho.Options{}))

	app.Logger.Fatal(app.Start(":10000"))
}
