package handlers

import (
	dbmodels "effectivemobiletest/internal/models/db-models"
	transportmodels "effectivemobiletest/internal/models/transport-models"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/render"
)

type UserGetter interface {
	ReadUsers(filters map[string]string, limit int, offset int) ([]dbmodels.UserOut, error)
}

func NewGetter(log *slog.Logger, ud UserGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()
		outValues := make(map[string]string)
		var err error
		offset, limit := 0, 0
		for k := range values {
			if k == "offset" {
				offset, err = strconv.Atoi(values.Get(k))

				if err != nil {
					log.Error("Wrong url params")
					render.JSON(w, r, transportmodels.Response{
						Status: http.StatusBadRequest,
						Error:  "Wrong url params"})
					return
				}
				continue
			} else if k == "limit" {
				limit, err = strconv.Atoi(values.Get(k))

				if err != nil {
					log.Error("Wrong url params")
					render.JSON(w, r, transportmodels.Response{
						Status: http.StatusBadRequest,
						Error:  "Wrong url params"})
					return
				}
				continue
			}

			outValues[k] = values.Get(k)
		}

		res, err := ud.ReadUsers(outValues, limit, offset)

		if err != nil {
			log.Error("Get users error", slog.String("error", err.Error()))
			render.JSON(w, r, transportmodels.Response{
				Status: http.StatusInternalServerError,
				Error:  "Get users error " + err.Error()})
			return
		}

		log.Info("Users getting")
		render.JSON(w, r, transportmodels.Response{
			Status: http.StatusOK,
			Users:  res,
		})
	}
}
