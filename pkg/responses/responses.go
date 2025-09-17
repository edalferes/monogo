package responses

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type CreatedResponse struct {
	ID string `json:"id"`
}

func OK(c echo.Context, response interface{}) error {
	return c.JSON(http.StatusOK, response)
}

func NoContent(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

func NotFound(c echo.Context, response interface{}) error {
	return c.JSON(http.StatusNotFound, response)
}

func BadRequest(c echo.Context, err error) error {
	return c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
}

func InternalServerError(c echo.Context, err error) error {
	return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
}

func Created(c echo.Context, id string) error {
	location := fmt.Sprintf("%s/%s", c.Request().URL.Path, id)
	c.Response().Header().Set("Location", location)
	return c.JSON(http.StatusCreated, CreatedResponse{ID: id})
}
