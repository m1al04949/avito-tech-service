package adduser

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
	UserID int `json:"user_id"`
	Method string
}

type UserSaver interface {
	SaveUser(int) error
}

func AddUser(log *slog.Logger, userSaver UserSaver) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.adduser"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

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

		err = userSaver.SaveUser(user)
		if errors.Is(err, storage.ErrUserExists) {
			log.Info("user already exists", slog.Int("user", user))
		}
		if err != nil {
			log.Error("failed to save user", logger.Err(err))
			render.JSON(w, r, response.Error("failed to save user"))
			return
		}

		log.Info("user added", slog.Int("user", user))

		render.JSON(w, r, Response{
			Response: response.OK(),
			UserID:   user,
			Method:   r.Method,
		})
	}
}
