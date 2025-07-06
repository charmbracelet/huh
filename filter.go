package huh

import "strings"

type FilterFunc func(target string, filter string) bool

func defaultFilterFunc(target string, filter string) bool {
	return strings.Contains(strings.ToLower(target), strings.ToLower(filter))
}
