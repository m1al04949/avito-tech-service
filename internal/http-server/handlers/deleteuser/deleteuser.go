package deleteuser

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
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

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", logger.Err(err))

			render.JSON(w, r, response.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", logger.Err(err))

			render.JSON(w, r, response.ValidationError(validateErr))

			return
		}

		user := req.UserID
		err = userDeleter.DeleteUser(user)

		if errors.Is(err, storage.ErrUserNotExists) {
			log.Info("user not exists", slog.Int("user", user))

			render.JSON(w, r, response.Error("user not exists"))

			return
		}
		if errors.Is(err, storage.ErrUserDelete) {
			log.Info("delete user from user segments table", slog.Int("user", user))

			render.JSON(w, r, response.Error("delete user from user segments table"))

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
