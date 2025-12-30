package utils

func MergeMap(source map[string][]string, target map[string][]string) {
	for k, v := range target {
		source[k] = v
	}
}
