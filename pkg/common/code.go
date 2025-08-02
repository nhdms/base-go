package common

const (
	AuthCodeNoEndPoint = iota + 100
	AuthCodeNoToken
	AuthCodeInvalidToken
	AuthCodeExpiredToken
	AuthCodeInvalidUser
	AuthCodeUserNotVerified
	AuthCodeUserBlocked
	AuthCodeUnauthorized
)

var Code2Message = map[int]string{
	AuthCodeNoEndPoint:      "No endpoint provided",
	AuthCodeNoToken:         "No token provided",
	AuthCodeInvalidToken:    "Invalid token",
	AuthCodeExpiredToken:    "Expired token",
	AuthCodeInvalidUser:     "Invalid user",
	AuthCodeUserNotVerified: "User not verified",
	AuthCodeUserBlocked:     "User blocked",
	AuthCodeUnauthorized:    "Unauthorized",
}
