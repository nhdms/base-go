package common

import (
	"database/sql"
	"errors"
	"strings"
)

var (
	ErrorServerError        = errors.New("server error")
	UnauthorizedError       = errors.New("unauthorized")
	ErrorUserInactivated    = errors.New("inactivated user")
	ErrorCompanyInactivated = errors.New("inactivated company")
	UserHasNoProfile        = errors.New("user has no profile")
	SQLNotFoundError        = errors.New("grpc not found. sql.ErrNoRows") // create new error because cannot compare error by grpc protocol
)

func IsNotFoundError(err error) bool {
	if err == nil {
		return false
	}

	return errors.Is(err, sql.ErrNoRows) || strings.Contains(err.Error(), SQLNotFoundError.Error())
}
