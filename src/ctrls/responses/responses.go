package responses

import (
	"net/http"

	"github.com/wakuwaku3/account-book.api/src/domains/apperrors"

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
func WriteErrorResponse(c echo.Context, err error) error {
	if cErr, ok := err.(apperrors.ClientError); ok {
		return c.JSON(http.StatusBadRequest, ErrorResponse{
			Errors: *cErr.GetErrorCodes(),
		})
	}
	c.Logger().Error(err, c.Request())
	return err
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
