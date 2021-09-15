package parse

func IsReserved(s string) bool {
	if s == "_" || s == "foundation" {
		return true
	}

	return false
}
