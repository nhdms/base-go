package dbtool

func GetQueryLimit(limit, defaultLimit, max uint64) uint64 {
	if limit <= 0 {
		return defaultLimit
	}

	if limit > max {
		return max
	}

	return limit
}

func GetQueryOffset(limit, page uint64) uint64 {
	if page < 1 {
		page = 1
	}

	return (page - 1) * limit
}
