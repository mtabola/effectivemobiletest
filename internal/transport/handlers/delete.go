package handlers

import (
	transportmodels "effectivemobiletest/internal/models/transport-models"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type UserDeleteter interface {
	DeleteUser(userId int) error
}

func NewDeleteter(log *slog.Logger, ud UserDeleteter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.delete.NewDeleteter"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		uid := r.URL.Query().Get("user_id")
		if uid == "" {
			log.Error("Wrong url params")
			render.JSON(w, r, transportmodels.Response{
				Status: http.StatusBadRequest,
				Error:  "Wrong url params"})
			return
		}

		uidNum, err := strconv.Atoi(uid)
		if err != nil {
			log.Error("UserId conversion failed", slog.String("error", err.Error()))
			render.JSON(w, r, transportmodels.Response{
				Status: http.StatusBadRequest,
				Error:  "UserId conversion failed! Please check your paramss"})
			return
		}

		err = ud.DeleteUser(uidNum)
		if err != nil {
			log.Error("User deletion failed", slog.String("error", err.Error()))
			render.JSON(w, r, transportmodels.Response{
				Status: http.StatusInternalServerError,
				Error:  "UserId deletion failed"})
			return
		}

		log.Info("User deleted")
		render.JSON(w, r, transportmodels.Response{
			Status:  http.StatusInternalServerError,
			Message: "User succssfuly deleted"})
	}
}
