package deleteuser

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/m1al04949/avito-tech-service/internal/lib/response"
	"github.com/m1al04949/avito-tech-service/internal/logger"
	"github.com/m1al04949/avito-tech-service/internal/storage"
	"golang.org/x/exp/slog"
)

type Request struct {
	UserID int `json:"user_id"`
}

type Response struct {
	response.Response
	UserID int `json:"user_id"`
	Method string
}

type UserDeleter interface {
	DeleteUser(int) error
}

func DeleteUser(log *slog.Logger, userDeleter UserDeleter) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.deleteuser"

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

		err = userDeleter.DeleteUser(user)

		if errors.Is(err, storage.ErrUserNotExists) {
			log.Info("user not exists", slog.Int("user", user))

			render.JSON(w, r, response.Error("user not exists"))

			return
		}
		if err != nil {
			log.Error("failed to delete user", logger.Err(err))

			render.JSON(w, r, response.Error("failed to delete user"))

			return
		}

		log.Info("user deleted", slog.Int("user", user))

		render.JSON(w, r, Response{
			Response: response.OK(),
			UserID:   user,
			Method:   r.Method,
		})
	}
}
