package delete

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"link-shortener/internal/lib/api/response"
	"link-shortener/internal/lib/logger/sl"
	"link-shortener/internal/storage"
	"log/slog"
	"net/http"
	"strconv"
)

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLDeleter
type URLDeleter interface {
	DeleteURL(ID int64) error
}

func New(log *slog.Logger, deleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			log.Error("can't parse url id", sl.Err(err))
			render.JSON(w, r, response.Error("invalid id"))
			return
		}

		err = deleter.DeleteURL(id)

		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url id not found", slog.Int64("id", id))

			render.JSON(w, r, response.Error("url id not found"))

			return
		}

		if err != nil {
			log.Error("failed to delete url", sl.Err(err))

			render.JSON(w, r, response.Error("failed to delete url"))

			return
		}

		log.Info("url deleted", slog.Int64("id", id))

		render.JSON(w, r, response.OK())
	}
}
