package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type SurveyElement struct {
	ID         string `json:"se_id"`
	FD_ID      string `json:"fd_id"`
	Is_control bool   `json:"is_control"`
}

func GetNextElement(c echo.Context) error {
	var result = SurveyElement{ID: "1234", FD_ID: "11357491", Is_control: false}
	var err error = nil
	//set assignment in assignment table.
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, result)
}
