package util

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func Ok(c echo.Context, response any) error {
	return c.JSON(http.StatusOK, response)
}

func OkMessage(c echo.Context, message string) error {
	return c.JSON(http.StatusOK, map[string]string{"message": message})
}

func BadRequest(c echo.Context, message string) error {
	return c.JSON(http.StatusBadRequest, map[string]string{"error": message})
}

func NotFound(c echo.Context, message string) error {
	return c.JSON(http.StatusNotFound, map[string]string{"error": message})
}

func Created(c echo.Context, response any) error {
	return c.JSON(http.StatusCreated, response)
}

func InternalServerError(c echo.Context, message string) error {
	return c.JSON(http.StatusInternalServerError, map[string]string{"error": message})
}
