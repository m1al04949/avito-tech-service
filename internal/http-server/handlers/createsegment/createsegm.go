package createsegment

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
	Slug string `json:"slug"`
}

type Response struct {
	response.Response
	Segment string `json:"slug,omitempty"`
	Method  string
}

type SegmSaver interface {
	SaveSegm(segmToSave string) error
}

func NewSegment(log *slog.Logger, segmSaver SegmSaver) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.createsegment"

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
		if segment == "" {
			err := fmt.Errorf("segment name (slug) is empty")
			log.Error("segment name (slug) is empty", logger.Err(err))
			render.JSON(w, r, response.Error("segment name (slug) is empty"))
			return
		}

		err = segmSaver.SaveSegm(segment)

		if errors.Is(err, storage.ErrSegmentExists) {
			log.Info("segment already exists", slog.String("segment", req.Slug))

			render.JSON(w, r, response.Error("segment already exists"))

			return
		}
		if err != nil {
			log.Error("failed to add segment", logger.Err(err))

			render.JSON(w, r, response.Error("failed to add segment"))

			return
		}

		log.Info("segment added", slog.String("segment", segment))

		render.JSON(w, r, Response{
			Response: response.OK(),
			Segment:  segment,
			Method:   r.Method,
		})
	}
}
