package util

func MergeMaps[K comparable, V interface{}](maps ...map[K]V) map[K]V {
	ret := make(map[K]V)

	for _, m := range maps {
		for k, v := range m {
			ret[k] = v
		}
	}

	return ret
}
