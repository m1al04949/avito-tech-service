package addtouser

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/m1al04949/avito-tech-service/internal/lib/response"
	"github.com/m1al04949/avito-tech-service/internal/lib/response/segmentsconv"
	"github.com/m1al04949/avito-tech-service/internal/logger"
	"github.com/m1al04949/avito-tech-service/internal/model"
	"github.com/m1al04949/avito-tech-service/internal/storage"
	"golang.org/x/exp/slog"
)

type Request struct {
	Segments []model.Segment `json:"segments,omitempty"`
}

type Response struct {
	response.Response
	UserID   int             `json:"user_id"`
	Segments []model.Segment `json:"segments"`
	Method   string
}

type UserSegmSaver interface {
	SaveSegmToUser(int, []string) error
}

func AddToUser(log *slog.Logger, userSegmSaver UserSegmSaver) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.addtouser"

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

		var req Request

		err = render.DecodeJSON(r.Body, &req)
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

		segms := req.Segments
		segments := segmentsconv.SegmentsConv(segms)

		err = userSegmSaver.SaveSegmToUser(user, segments)
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
			Segments: segms,
			Method:   r.Method,
		})
	}

}
