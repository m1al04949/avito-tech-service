package deletefromuser

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

type UserSegmDeleter interface {
	DeleteSegmFromUser(int, []string) error
}

func DeleteFromUser(log *slog.Logger, userSegmDeleter UserSegmDeleter) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.deletefromuser"

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

		err = userSegmDeleter.DeleteSegmFromUser(user, segments)
		if errors.Is(err, storage.ErrSegmentsNotExists) {
			for _, v := range segments {
				log.Info("segment not exists", slog.String("segment", v))
			}
			render.JSON(w, r, response.Error("segments not exists"))
			return
		}
		if errors.Is(err, storage.ErrUserNotExists) {
			log.Error("user not exists", logger.Err(err))
			render.JSON(w, r, response.Error("user not exists"))
			return
		}
		if err != nil {
			log.Error("failed to delete segments from user", logger.Err(err))
			render.JSON(w, r, response.Error("failed to delete segments from user"))
			return
		}

		log.Info("segments deleted from user", slog.Int("user", user))

		render.JSON(w, r, Response{
			Response: response.OK(),
			UserID:   user,
			Segments: segms,
			Method:   r.Method,
		})
	}

}
