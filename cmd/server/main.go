package main

import (
	"fmt"
	"github.com/acai-travel/tech-challenge/internal/observability"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"log"
	"log/slog"
	"net/http"

	"context"
	"github.com/acai-travel/tech-challenge/internal/chat"
	"github.com/acai-travel/tech-challenge/internal/chat/assistant"
	"github.com/acai-travel/tech-challenge/internal/chat/model"
	"github.com/acai-travel/tech-challenge/internal/httpx"
	"github.com/acai-travel/tech-challenge/internal/mongox"
	"github.com/acai-travel/tech-challenge/internal/pb"
	"github.com/gorilla/mux"
	"github.com/twitchtv/twirp"
)

func main() {
	shutdownTracer, err := observability.InitTracer()
	if err != nil {
		log.Fatal(err)
	}
	defer shutdownTracer(context.Background())

	shutdownMetrics := observability.InitMetrics(context.Background())
	defer shutdownMetrics()

	//Exposing metrics prometheus
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Fatal(http.ListenAndServe(":2112", nil))
	}()

	mongo := mongox.MustConnect()

	repo := model.New(mongo)
	assist := assistant.New()

	server := chat.NewServer(repo, assist)

	// Configure handler
	handler := mux.NewRouter()
	handler.Use(
		httpx.Logger(),
		httpx.Recovery(),
	)

	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, "Hi, my name is Clippy!")
	})

	handler.PathPrefix("/twirp/").Handler(pb.NewChatServiceServer(server, twirp.WithServerJSONSkipDefaults(true)))

	// Start the server
	slog.Info("Starting the server...")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		panic(err)
	}
}
