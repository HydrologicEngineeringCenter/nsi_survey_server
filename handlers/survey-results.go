package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type SurveyResult struct {
	ID    string `json:"sr_id"`
	FD_ID string `json:"fd_id"`
}

func PostSurveyResult(c echo.Context) error {
	var result = SurveyResult{ID: "1234", FD_ID: "11357491"}
	var err error = nil
	//set assignment in assignment table.
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, result)
}
