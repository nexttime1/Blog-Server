package image

func InList(key string, list []string) bool {
	for _, value := range list {
		if key == value {
			return true
		}
	}
	return false
}
