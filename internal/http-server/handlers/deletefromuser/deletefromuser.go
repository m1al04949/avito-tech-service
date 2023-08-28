package deletefromuser

import (
	"net/http"

	"github.com/m1al04949/avito-tech-service/internal/lib/response"
	"golang.org/x/exp/slog"
)

type Segment struct {
	Slug string `json:"slug,omitempty"`
}

type Request struct {
	UserID   int       `json:"user_id"`
	Segments []Segment `json:"segments,omitempty"`
}

type Response struct {
	response.Response
	UserID   int       `json:"user_id"`
	Segments []Segment `json:"segments,omitempty"`
	Method   string
}

type UserSegmDeleter interface {
	DeleteSegmFromUser(int, []string) error
}

func DeleteFromUser(log *slog.Logger, userSegmDeleter UserSegmDeleter) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		// const op = "handlers.adduser"

		// log = log.With(
		// 	slog.String("op", op),
		// 	slog.String("request_id", middleware.GetReqID(r.Context())),
		// )

		// var req Request

		// err := render.DecodeJSON(r.Body, &req)
		// if err != nil {
		// 	log.Error("failed to decode request body", logger.Err(err))

		// 	render.JSON(w, r, response.Error("failed to decode request"))

		// 	return
		// }

		// log.Info("request body decoded", slog.Any("request", req))

		// if err := validator.New().Struct(req); err != nil {
		// 	validateErr := err.(validator.ValidationErrors)

		// 	log.Error("invalid request", logger.Err(err))

		// 	render.JSON(w, r, response.ValidationError(validateErr))

		// 	return
		// }

		// user := req.UserID

		// err = userSaver.SaveUser(user)
		// if errors.Is(err, storage.ErrUserExists) {
		// 	log.Info("user already exists", slog.Int("user", req.UserID))
		// }
		// if err != nil {
		// 	log.Error("failed to save user", logger.Err(err))
		// 	render.JSON(w, r, response.Error("failed to save user"))
		// 	return
		// }

		// log.Info("user added", slog.Int("user", user))

		// render.JSON(w, r, Response{
		// 	Response: response.OK(),
		// 	UserID:   user,
		// 	Method:   r.Method,
		// })
	}

}
