package main

import (
	"context"
	"net/http"

	"github.com/google/go-github/v57/github"
	"github.com/labstack/echo/v4"
)

func main() {
	app := echo.New()
	client := github.NewClient(nil)
	repoOptions := &github.RepositoryListByUserOptions{Type: "public"}

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

	app.Logger.Fatal(app.Start(":10000"))
}