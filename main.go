package main

import (
	"context"
	"embed"
	"net/http"
	"x-straight-check/config"
	"x-straight-check/pkg/gemini"
	"x-straight-check/pkg/xscraper"
	"x-straight-check/views"
)

//go:embed assets/*
var assets embed.FS

func main() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "version", config.Env().Version)

	http.Handle("/assets/", http.FileServer(http.FS(assets)))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		views.Layout(views.TheThing()).Render(ctx, w)
	})

	http.HandleFunc("/api/check", func(w http.ResponseWriter, r *http.Request) {
		xHandle := r.FormValue("x_handle")
		if xHandle == "" {
			views.Error("Empty X Handle!").Render(ctx, w)
			return
		}

		posts, user, err := xscraper.GetUserPosts(xHandle, "en", 35)
		if err != nil {
			panic(err)
		}

		res, err := gemini.CheckUserStraightness(posts, user)
		if err != nil {
			panic(err)
		}

		views.AnalysisResult(*res).Render(ctx, w)
	})

	http.ListenAndServe(":"+config.Env().Port, nil)

}
