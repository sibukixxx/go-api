package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

// Handler は HTTP ハンドラーの構造体です
type Handler struct{}

// NewHandler は Handler のインスタンスを生成します
func NewHandler() *Handler {
	return &Handler{}
}

// Hello は簡単な Hello World エンドポイントです
func (h *Handler) Hello(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"message": "Hello, World!"})
}

// Router は mux.Router のインスタンスを生成し、ルートを設定します
func NewRouter(handler *Handler) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/hello", handler.Hello).Methods("GET")
	return r
}

// Server は http.Server のインスタンスを生成します
func NewServer(lc fx.Lifecycle, router *mux.Router) *http.Server {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Println("Starting HTTP server on :8080")
			go srv.ListenAndServe()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Stopping HTTP server")
			return srv.Shutdown(ctx)
		},
	})

	return srv
}

func main() {
	fx.New(
		fx.Provide(
			NewHandler,
			NewRouter,
			NewServer,
		),
		fx.Invoke(func(*http.Server) {}),
	).Run()
}
