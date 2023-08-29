package getuser

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/m1al04949/avito-tech-service/internal/lib/response"
	"github.com/m1al04949/avito-tech-service/internal/logger"
	"github.com/m1al04949/avito-tech-service/internal/storage"
	"golang.org/x/exp/slog"
)

type Response struct {
	response.Response
	UserID   int      `json:"user_id"`
	Segments []string `json:"segments"`
	Method   string
}

type UserGetter interface {
	GetUser(int) ([]string, error)
}

func GetFromUser(log *slog.Logger, userGetter UserGetter) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.getuser"

		id := chi.URLParam(r, "id")
		if id == "" {
			log.Info("id is empty")

			render.JSON(w, r, response.Error("invalid request"))

			return
		}
		user, err := strconv.Atoi(id)
		if err != nil {
			log.Info("id is not int")

			render.JSON(w, r, response.Error("invalid id"))

			return
		}

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		segments, err := userGetter.GetUser(user)
		if errors.Is(err, storage.ErrUserNotExists) {
			log.Info("user not exists", slog.Int("user", user))
		}
		if err != nil {
			log.Error("failed to get user info", logger.Err(err))
			render.JSON(w, r, response.Error("failed to get user info"))
			return
		}

		log.Info("user info is getted", slog.Int("user", user))

		render.JSON(w, r, Response{
			Response: response.OK(),
			UserID:   user,
			Segments: segments,
			Method:   r.Method,
		})
	}
}
