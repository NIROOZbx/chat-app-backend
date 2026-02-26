package response

import "net/http"

const (
	StatusOK       = http.StatusOK
	StatusCreated  = http.StatusCreated
	StatusAccepted = http.StatusAccepted

	StatusBadRequest      = http.StatusBadRequest
	StatusUnauthorized    = http.StatusUnauthorized
	StatusForbidden       = http.StatusForbidden
	StatusNotFound        = http.StatusNotFound
	StatusConflict        = http.StatusConflict
	StatusTooManyRequests = http.StatusTooManyRequests

	StatusInternalServerError = http.StatusInternalServerError
	StatusNotImplemented      = http.StatusNotImplemented
	StatusBadGateway          = http.StatusBadGateway
	StatusServiceUnavailable  = http.StatusServiceUnavailable
	StatusGatewayTimeout      = http.StatusGatewayTimeout
)
const (
    SuccessMsgCreated = "created successfully"
    SuccessMsgUpdated = "updated successfully"
    SuccessMsgDeleted = "deleted successfully"
    SuccessMsgFetched = "Data fetched successfully"
)

const (
	MaxMember="maximum 50 members allowed"
	UserNotMember="user is not member of this group"
	UserNotFound="user was not found"
)