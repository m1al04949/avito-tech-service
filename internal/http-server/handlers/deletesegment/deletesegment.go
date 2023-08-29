package deletesegment

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
	Slug string `json:"slug"`
}

type Response struct {
	response.Response
	Segment string `json:"slug"`
	Method  string
}

type SegmDeleter interface {
	DeleteSegm(string) error
}

func DelSegment(log *slog.Logger, segmDeleter SegmDeleter) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.deletesegment"

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

		segment := req.Slug

		err = segmDeleter.DeleteSegm(segment)

		if errors.Is(err, storage.ErrSegmentNotExists) {
			log.Info("segment not exists", slog.String("segment", req.Slug))

			render.JSON(w, r, response.Error("segment not exists"))

			return
		}
		if errors.Is(err, storage.ErrSegmentDelete) {
			log.Info("delete segment from user segments table", slog.String("segment", segment))

			render.JSON(w, r, response.Error("delete segment from user segments table"))

			return
		}
		if err != nil {
			log.Error("failed to delete segment", logger.Err(err))

			render.JSON(w, r, response.Error("failed to delete segment"))

			return
		}

		log.Info("segment deleted", slog.String("segment", segment))

		render.JSON(w, r, Response{
			Response: response.OK(),
			Segment:  segment,
			Method:   r.Method,
		})
	}
}
