package adduser

import (
	"errors"
	"fmt"
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

//go:generate go run github.com/vektra/mockery/v2 --name=UserSaver
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
		if user == 0 {
			err := fmt.Errorf("user_id is empty")
			log.Error("user_id is empty", logger.Err(err))
			render.JSON(w, r, response.Error("user_id is empty"))
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
