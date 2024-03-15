package main

import (
	"net/http"
)

type HealthCheckResponse struct {
	Status     string `json:"status"`
	SystemInfo struct {
		Environment string `json:"environment"`
	} `json:"system_info"`
}

// @Summary Healthcheck
// @Description Check the health status of the application
// @Tags Healthcheck
// @Accept json
// @Produce json
// @Success 200 {object} HealthCheckResponse
// @Router /healthcheck [get]
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
		},
	}

	err := app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
