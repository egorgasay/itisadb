package pkg

func IsTheSameArray[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}

	tmp := make(map[T]struct{})
	for _, el := range a {
		tmp[el] = struct{}{}
	}
	for _, el := range b {
		if _, ok := tmp[el]; !ok {
			return false
		}
	}
	return true
}
