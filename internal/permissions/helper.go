package permissions

func HasPermission(userPermissions int64, requiredPermission int64) bool {
	return userPermissions&requiredPermission != 0
}
