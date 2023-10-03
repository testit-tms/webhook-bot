package send

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/testit-tms/webhook-bot/internal/entities"
	"github.com/testit-tms/webhook-bot/internal/lib/handlers"
	"github.com/testit-tms/webhook-bot/internal/lib/logger/sl"
	val "github.com/testit-tms/webhook-bot/internal/lib/validator"
	"golang.org/x/exp/slog"
)

//go:generate mockgen -source=$GOFILE -destination=$PWD/mocks/${GOFILE} -package=mocks
type sender interface {
	SendMessage(ctx context.Context, msg entities.Message) error
}

// New returns a new http.HandlerFunc that sends a message using the provided sender.
// It validates the request, converts it to a message, and sends it using the sender.
// If any error occurs during the process, it returns an error response.
// It requires an Authorization token in the request header.
func New(log *slog.Logger, sender sender) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "transport.rest.send.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		token := r.Header.Get("Authorization")
		if token == "" {
			log.Debug("token not found")
			handlers.NewErrorResponse(w, http.StatusUnauthorized, "token is required")
			return
		}

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			handlers.NewErrorResponse(w, http.StatusBadRequest, "failed to decode request")
			return
		}

		log.Debug("request body decoded", slog.Any("request", req))

		v := validator.New()
		// nolint:errcheck
		v.RegisterValidation("parse-mode", val.ValidateParseMode)
		if err := v.Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(err))

			handlers.NewErrorResponse(w, http.StatusBadRequest, handlers.ValidationError(validateErr))

			return
		}

		message := req.convertToDomain()
		message.Token = token

		log.Debug("request convert to message", slog.Any("message", message))
		err = sender.SendMessage(r.Context(), message)
		if err != nil {
			log.Error("can not send message", sl.Err(err))

			handlers.NewErrorResponse(w, http.StatusInternalServerError, "can't send message")
			return
		}

		w.WriteHeader(http.StatusOK)
		// nolint:errcheck
		w.Write([]byte("message sent"))
	}
}
