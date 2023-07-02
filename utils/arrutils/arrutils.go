package arrutils

import (
	"sort"
)

/// возвращает массив уникальных строк.
func UniqueStr(arr []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range arr {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// Sort возвращает отсортированный массив строк.
func Sort(arr []string) []string {
	sort.Strings(arr)
	return arr
}
