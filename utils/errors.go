package utils

import (
	"errors"
	"go-chi-boilerplate/utils/constants"
	"net/http"
)

// ERROR HTTP
var (
	// 5XX
	ErrorInternalServer   = errors.New(constants.INTERNAL_SERVER_ERROR)
	ErrorStartTransaction = errors.New(constants.TRANSACTION_FAILED)

	// 4XX
	ErrorUnauthorized        = errors.New(constants.UNAUTHORIZED)
	ErrorBadRequest          = errors.New(constants.BAD_REQUEST)
	ErrorNotFound            = errors.New("your requested item is not found")
	ErrorDuplicateData       = errors.New("your item already exist")
	ErrorMaxSize             = errors.New("maximum size exceeded")
	ErrorInvalidBearerToken  = errors.New("unauthorized: Invalid Bearer token")
	ErrorLoginRequired       = errors.New("unauthorized: " + constants.LOGIN_REQUIRED)
	ErrorBearer              = errors.New("unauthorized: Bearer token is missing or empty")
	ErrorState               = errors.New("unauthorized: State is invalid")
	ErrorRefreshTokenRevoked = errors.New("refresh token revoked")
	ErrorAccessTokenExpired  = errors.New("access token is expired")
	ErrorInvalidToken        = errors.New(constants.INVALID_TOKEN)

	// 2XX
	ErrorNoContent = errors.New(constants.NO_CONTENT)

	// other
	ErrorBuildQuery = errors.New("error building query")
)

func GetStatusCode(err error) (int, string) {
	if err == nil {
		return http.StatusOK, constants.SUCCESS
	}

	switch err {
	case ErrorBadRequest, ErrorInvalidToken:
		return http.StatusBadRequest, constants.BAD_REQUEST
	case ErrorNotFound:
		return http.StatusNotFound, constants.DATA_NOT_EXIST
	case ErrorDuplicateData:
		return http.StatusConflict, constants.ALREADY_EXISTS
	case ErrorUnauthorized, ErrorState, ErrorBearer, ErrorInvalidBearerToken, ErrorLoginRequired:
		return http.StatusUnauthorized, constants.UNAUTHORIZED
	case ErrorRefreshTokenRevoked:
		return http.StatusUnauthorized, constants.REFRESH_TOKEN_REVOKED
	case ErrorNoContent:
		return http.StatusNoContent, constants.NO_CONTENT
	case ErrorAccessTokenExpired:
		return http.StatusUnauthorized, constants.ACCESS_TOKEN_EXPIRED
	case ErrorMaxSize:
		return http.StatusBadRequest, constants.MAX_SIZE_EXCEEDED
	default:
		return http.StatusInternalServerError, constants.INTERNAL_SERVER_ERROR
	}
}

func ErrorQueryBuilder(tableName string, err error) string {
	result := ErrorBuildQuery.Error() + " " + tableName + " : " + err.Error()
	return result
}
