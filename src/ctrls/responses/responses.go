package responses

import (
	"net/http"

	"github.com/labstack/echo"
)

type (
	// ErrorResponse はエラー時のレスポンスです
	ErrorResponse struct {
		Errors []string `json:"errors"`
	}
	// Response はレスポンスです
	Response struct {
		Result interface{} `json:"result"`
	}
)

// WriteErrorResponse はエラーをレスポンスボディーに書き込みます
func WriteErrorResponse(c echo.Context, errs ...error) error {
	errors := make([]string, len(errs))
	for i, err := range errs {
		errors[i] = err.Error()
	}
	return c.JSON(http.StatusBadRequest, ErrorResponse{
		Errors: errors,
	})
}

// WriteResponse は結果をレスポンスボディーに書き込みます
func WriteResponse(c echo.Context, result interface{}) error {
	return c.JSON(http.StatusOK, Response{
		Result: result,
	})
}

// WriteEmptyResponse は空の結果をレスポンスボディーに書き込みます
func WriteEmptyResponse(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

// WriteUnAuthorizedErrorResponse は空の結果をレスポンスボディーに書き込みます
func WriteUnAuthorizedErrorResponse(c echo.Context) error {
	return c.NoContent(http.StatusUnauthorized)
}
