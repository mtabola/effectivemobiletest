package handlers

import (
	"effectivemobiletest/internal/logger"
	dbmodels "effectivemobiletest/internal/models/db-models"
	enrichmodels "effectivemobiletest/internal/models/enrich-models"
	transportmodels "effectivemobiletest/internal/models/transport-models"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

const (
	AgifyLink  = "https://api.agify.io/?name="
	GenderLink = "https://api.genderize.io/?name="
	RegionLink = "https://api.nationalize.io/?name="
)

type UserSaver interface {
	CreateUser(user dbmodels.UserIn) error
	RegionCodeChekingAndCreating(code string) (int, error)
	GetGenderId(gn string) (int, error)
	GetRandomRegionId() (int, error)
}

func NewSaver(log *slog.Logger, us UserSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.save.NewSaver"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req transportmodels.Request

		log.Debug("Request: ", r, r.Body)

		err := render.Decode(r, &req)
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
		log.Debug("Message decoded:", slog.Any("request", req))

		var age enrichmodels.Agify
		var gender enrichmodels.Genderize
		var region enrichmodels.Nationalize
		enrichRequest(AgifyLink+req.FName, &age)
		enrichRequest(GenderLink+req.FName, &gender)
		enrichRequest(RegionLink+req.FName, &region)

		user := dbmodels.UserIn{
			FName: req.FName,
			SName: req.SName,
			PName: req.PName,
			Age:   uint8(age.Age),
		}

		var reg int
		if len(region.Countries) == 0 {
			reg, err = us.GetRandomRegionId()
			if err != nil {
				log.Error("Country unmarshaling error", slog.String("error", err.Error()))
				render.JSON(w, r, transportmodels.Response{
					Status: http.StatusInternalServerError,
					Error:  fmt.Sprintf("Country decoding error: %v", err)})
				return
			}
		} else {
			reg, err = us.RegionCodeChekingAndCreating(region.Countries[0].CountryId)
			if err != nil {
				log.Error("Country unmarshaling error", slog.String("error", err.Error()))
				render.JSON(w, r, transportmodels.Response{
					Status: http.StatusInternalServerError,
					Error:  fmt.Sprintf("Country decoding error: %v", err)})
				return
			}
		}
		user.RegionId = reg

		gdr, err := us.GetGenderId(gender.Gender)
		if err != nil {
			log.Error("Gender unmarshaling error", slog.String("error", err.Error()))
			render.JSON(w, r, transportmodels.Response{
				Status: http.StatusInternalServerError,
				Error:  fmt.Sprintf("Gender decoding error: %v", err)})
			return
		}
		user.GenderId = gdr

		err = us.CreateUser(user)
		if err != nil {
			log.Error("User insertion error", slog.String("error", err.Error()))
			render.JSON(w, r, transportmodels.Response{
				Status: http.StatusInternalServerError,
				Error:  fmt.Sprintf("User insertion error: %v", err)})
			return
		}

		log.Info("User added")
		render.JSON(w, r, transportmodels.Response{
			Status:  http.StatusOK,
			Message: "User added successfully",
		})
	}
}

func enrichRequest(address string, data interface{}) error {
	r, err := http.Get(address)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)

	err = decoder.Decode(data)
	if err != nil {
		return err
	}
	return nil
}
