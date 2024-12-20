package main

import (
	"context"
	"embed"
	"net/http"
	"x-straight-check/config"
	"x-straight-check/pkg/ratelimiter"
	"x-straight-check/views"
)

//go:embed assets/*
var assets embed.FS

func main() {
	ctx := context.Background()
	version := config.Env().Version
	if len(version) >= 7 {
		version = version[:7]
	}
	ctx = context.WithValue(ctx, "version", version)

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
		res, err := ratelimiter.GetUserAnalysis(xHandle, "en")
		if err != nil {
			views.Error("Internal server error").Render(ctx, w)
			return
		}

		views.AnalysisResult(*res).Render(ctx, w)
	})

	http.ListenAndServe(":"+config.Env().Port, nil)

}
