package save

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"io"
	"link-shortener/internal/lib/api/response"
	"link-shortener/internal/lib/logger/sl"
	"link-shortener/internal/lib/random"
	"link-shortener/internal/storage"
	"log/slog"
	"net/http"
	"regexp"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
	ID    int64  `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

// ref to conf
const aliasLength = 6

type URLSaver interface {
	SaveURL(URL string, alias string) (int64, error)
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")
			render.JSON(w, r, response.Error("empty request"))
			return
		}

		if err != nil {
			log.Error("failed to parse request body", sl.Err(err))
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			var validateErr validator.ValidationErrors
			errors.As(err, &validateErr)
			log.Error("failed to validate request", sl.Err(err))
			render.JSON(w, r, response.ValidationError(validateErr))
			return
		}

		// Alias validation: ensure no special characters
		if req.Alias != "" && !isValidAlias(req.Alias) {
			log.Info("invalid alias", slog.String("alias", req.Alias))
			render.JSON(w, r, response.Error("invalid alias (special characters not allowed)"))
			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}

		id, err := urlSaver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrURLExist) {
			log.Info("url already exists", slog.String("url", req.URL))
			render.JSON(w, r, response.Error("url already exists"))
			return
		}
		log.Info("url saved", slog.Int64("id", id))
		responseOK(w, r, alias, id)
	}
}

// Sends a successful response with a custom JSON payload.
func responseOK(w http.ResponseWriter, r *http.Request, alias string, id int64) {
	render.JSON(w, r, Response{
		Response: response.OK(),
		Alias:    alias,
		ID:       id,
	})
}

// Custom validation for the alias to only allow alphanumeric characters, hyphens, and underscores
func isValidAlias(alias string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	return re.MatchString(alias)
}
