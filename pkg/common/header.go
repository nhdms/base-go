package common

const (
	HeaderUserId       = "X-AT-UserId"
	HeaderChildUserIds = "X-AT-Child-UserIds"
)

var ExtraDataHeaders = []string{
	HeaderUserId,
	HeaderChildUserIds,
	//"X-AT-Shop-Id",
	//"X-AT-Currency",
	//"X-AT-Language",
	//"X-AT-Project-Id",
	//"X-AT-Version",
	// ... define header to pass for service here
	// it will be extracted from jwt token
}
