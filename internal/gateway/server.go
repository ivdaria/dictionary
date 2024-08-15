package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"dictionary/internal/config"
	"dictionary/internal/convert"
	"dictionary/internal/entity"
	er "dictionary/internal/errors"
	"dictionary/pkg/gateway/model"
)

type itemsRepo interface {
	GetItemByID(ctx context.Context, id int64) (*entity.TranslationItem, error)
	CreateItem(ctx context.Context, item *entity.TranslationItem) (int64, error)
	UpdateItem(ctx context.Context, item *entity.TranslationItem) error
	ListItems(ctx context.Context) ([]*entity.TranslationItem, error)
	DeleteItem(ctx context.Context, id int64) error
}

type AppServer struct {
	server *http.Server
	repo   itemsRepo
}

func NewAppServer(cfg *config.Config, repo itemsRepo) *AppServer {
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
	mux.HandleFunc("POST /items/{id}/edit", appServer.UpdateItem)
	mux.HandleFunc("GET /items/{id}", appServer.GetItemByID)
	mux.HandleFunc("GET /items", appServer.ListItems)
	mux.HandleFunc("DELETE /items/{id}", appServer.DeleteItem)

	server.Handler = appServer.corsMiddleware(server.Handler)
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

func (s *AppServer) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With, Accept, Origin, X-AUTH-SID, X-ACCESS-TOKEN")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS,HEAD")
		w.Header().Add("Access-Control-Allow-Origin", r.Header.Get("Origin"))

		// immediately response for preflight request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *AppServer) CreateItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// New decoder -чтобы десериализовать тело запроса в нашу модель
	// не обязательно джсон. можем прислать хмл и тд. Соответственно, тогда нужен будет другой декодер
	decoder := json.NewDecoder(r.Body)
	// создаем модель, в которую будем декодировать тело запроса
	mdl := model.CreateItemRequestBody{}
	if err := decoder.Decode(&mdl); err != nil {
		slog.ErrorContext(
			ctx,
			"create item",
			slog.String("error", fmt.Errorf("decode body to model: %w", err).Error()),
		)
		//
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	// получаем айтем из модели
	item := convert.ItemFromCreateItemRequestBody(&mdl)
	// вызываем нашу функцию
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

	// Модель на ответ создается для того, чтобы держать слой API отделенным от слоя бизнес логики
	responseMdl := model.CreateItemResponseBody{
		ID: id,
	}

	// сериализация модели на ответ
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

	// пишем заголовок о том, что слово успешно создано
	// нужно сначала заголовок писать, а потом тело.
	// Иначе сначала запишется тело, а потом заголовок сразу запишется в статус 200 ОК

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

func (s *AppServer) GetItemByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idString := r.PathValue("id")
	id, _ := strconv.ParseInt(idString, 10, 64)

	item, err := s.repo.GetItemByID(ctx, id)
	if err != nil {
		slog.ErrorContext(
			ctx,
			"get item by id",
			slog.String("error", fmt.Errorf("get item by repo: %w", err).Error()),
		)

		if errors.Is(err, er.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseMdl := model.GetItemByIDResponseBody{
		ID:          item.ID,
		Word:        item.Word,
		Translation: item.Translation,
	}

	responseMdlBytes, err := json.Marshal(responseMdl)
	if err != nil {
		slog.ErrorContext(
			ctx,
			"get item by id",
			slog.String("error", fmt.Errorf("marshall response: %w", err).Error()),
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(responseMdlBytes); err != nil {
		slog.ErrorContext(
			ctx,
			"get item by id",
			slog.String("error", fmt.Errorf("write response: %w", err).Error()),
		)
		return
	}
}

func (s *AppServer) UpdateItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idString := r.PathValue("id")

	decoder := json.NewDecoder(r.Body)
	mdl := model.UpdateItemRequestBody{}
	if err := decoder.Decode(&mdl); err != nil {
		slog.ErrorContext(
			ctx,
			"update item",
			slog.String("error", fmt.Errorf("decode body to model: %w", err).Error()),
		)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	id, _ := strconv.ParseInt(idString, 10, 64)
	item := convert.ItemFromUpdateItemRequestBody(id, &mdl)

	err := s.repo.UpdateItem(ctx, item)
	if err != nil {
		slog.ErrorContext(
			ctx,
			"update item",
			slog.String("error", fmt.Errorf("create item by repo: %w", err).Error()),
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *AppServer) ListItems(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	items, err := s.repo.ListItems(ctx)
	if err != nil {
		slog.ErrorContext(
			ctx,
			"list items",
			slog.String("error", fmt.Errorf("get items by repo: %w", err).Error()),
		)

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	responseMdl := model.ListItemsResponseBody{
		Items: make([]model.Item, 0, len(items)),
	}
	for _, item := range items {
		responseMdl.Items = append(responseMdl.Items, model.Item{
			ID:          item.ID,
			Word:        item.Word,
			Translation: item.Translation,
		})
	}

	responseMdlBytes, err := json.Marshal(responseMdl)
	if err != nil {
		slog.ErrorContext(
			ctx,
			"list items",
			slog.String("error", fmt.Errorf("marshall response: %w", err).Error()),
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(responseMdlBytes); err != nil {
		slog.ErrorContext(
			ctx,
			"list items",
			slog.String("error", fmt.Errorf("write response: %w", err).Error()),
		)
		return
	}
}

func (s *AppServer) DeleteItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idString := r.PathValue("id")
	id, _ := strconv.ParseInt(idString, 10, 64)

	err := s.repo.DeleteItem(ctx, id)
	if err != nil {
		slog.ErrorContext(
			ctx,
			"delete item",
			slog.String("error", fmt.Errorf("delete item: %w", err).Error()),
		)

		if errors.Is(err, er.ErrNoRowsAffected) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
