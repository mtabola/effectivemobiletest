package handlers

import (
	"effectivemobiletest/internal/logger"
	dbmodels "effectivemobiletest/internal/models/db-models"
	transportmodels "effectivemobiletest/internal/models/transport-models"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type UserUpdater interface {
	UpdateUser(user dbmodels.UserIn) error
	RegionCodeChekingAndCreating(code string) (int, error)
	GetGenderId(gender string) (int, error)
}

func NewUpdater(log *slog.Logger, uu UserUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.save.NewUpdater"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var uout dbmodels.UserOut

		err := render.Decode(r, &uout)
		if errors.Is(err, io.EOF) {
			log.Error("Request body is empty")

			render.JSON(w, r, transportmodels.Response{
				Status: http.StatusBadRequest,
				Error:  fmt.Sprintf("Request is empty: %v", err),
			})
			return
		} else if err != nil {
			log.Error("Failed to decode request body", logger.Err(err))

			render.JSON(w, r, transportmodels.Response{
				Status: http.StatusInternalServerError,
				Error:  fmt.Sprintf("Request decoding error: %v", err),
			})
			return
		}

		if uout.UserId == 0 || uout.Age == 0 || uout.FName == "" || uout.SName == "" || uout.GenderName == "" || uout.RegionCode == "" {
			render.JSON(w, r, transportmodels.Response{
				Status: http.StatusInternalServerError,
				Error:  "The following fields are required: UserId, FName, SName, Age, RegionCode, GenderName",
			})
			return
		}
		log.Debug("Message decoded:", slog.Any("request", uout))

		uin := dbmodels.UserIn{
			UserId: uout.UserId,
			FName:  uout.FName,
			SName:  uout.SName,
			PName:  uout.PName,
			Age:    uout.Age,
		}

		id, err := uu.RegionCodeChekingAndCreating(uout.RegionCode)
		if err != nil {
			log.Error("Getting region code error", logger.Err(err))

			render.JSON(w, r, transportmodels.Response{
				Status: http.StatusInternalServerError,
				Error:  fmt.Sprintf("Getting region code error: %v", err),
			})
			return
		}
		uin.RegionId = id

		id, err = uu.GetGenderId(uout.GenderName)
		if err != nil {
			log.Error("Getting gender id error", logger.Err(err))

			render.JSON(w, r, transportmodels.Response{
				Status: http.StatusInternalServerError,
				Error:  fmt.Sprintf("Getting gender id error: %v", err),
			})
			return
		}
		uin.GenderId = id

		err = uu.UpdateUser(uin)

		if err != nil {
			log.Error("User update error", logger.Err(err))

			render.JSON(w, r, transportmodels.Response{
				Status: http.StatusInternalServerError,
				Error:  fmt.Sprintf("User update error: %v", err),
			})
			return
		}

		log.Info("User updated")
		render.JSON(w, r, transportmodels.Response{
			Status:  http.StatusOK,
			Message: "User updated successfully",
		})
	}
}
