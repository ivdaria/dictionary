package gateway

import (
	"context"
	"dictionary/internal/config"
	"dictionary/internal/convert"
	translationitems "dictionary/internal/repository/translation-items"
	"dictionary/pkg/gateway/model"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type AppServer struct {
	server *http.Server
	repo   *translationitems.Repo
}

func NewAppServer(cfg *config.Config, repo *translationitems.Repo) *AppServer {
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    cfg.HTTP.ListenAddr,
		Handler: mux,
	}
	appServer := &AppServer{
		server: server,
		repo:   repo,
	}

	mux.HandleFunc("POST /items", appServer.CreateItem)

	return appServer
}

func (s *AppServer) Run() error {
	if err := s.server.ListenAndServe(); err != nil {
		return fmt.Errorf("run server: %w", err)
	}

	return nil
}

func (s *AppServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *AppServer) CreateItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	decoder := json.NewDecoder(r.Body)
	mdl := model.CreateItemRequestBody{}
	if err := decoder.Decode(&mdl); err != nil {
		slog.ErrorContext(
			ctx,
			"create item",
			slog.String("error", fmt.Errorf("decode body to model: %w", err).Error()),
		)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	item := convert.ItemFromCreateItemRequestBody(&mdl)
	id, err := s.repo.CreateItem(ctx, item)
	if err != nil {
		slog.ErrorContext(
			ctx,
			"create item",
			slog.String("error", fmt.Errorf("create item by repo: %w", err).Error()),
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseMdl := model.CreateItemResponseBody{
		ID: id,
	}

	responseMdlBytes, err := json.Marshal(responseMdl)
	if err != nil {
		slog.ErrorContext(
			ctx,
			"create item",
			slog.String("error", fmt.Errorf("marshall response: %w", err).Error()),
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(responseMdlBytes); err != nil {
		slog.ErrorContext(
			ctx,
			"create item",
			slog.String("error", fmt.Errorf("write response: %w", err).Error()),
		)
		return
	}
}
