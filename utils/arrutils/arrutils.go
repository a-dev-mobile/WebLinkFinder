package arrutils

import (
	"fmt"
	"sort"
)

// / возвращает массив уникальных строк.
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

func DeleteElement(slice []string, s string) []string {
	index := -1
	for i, item := range slice {
		if item == s {
			index = i
			break
		}
	}
	if index == -1 {
		fmt.Println("Element not found in the slice")
		return slice // Return original slice if element not found
	}
	// append the items before index with the items after index
	return append(slice[:index], slice[index+1:]...)
}

// Sort возвращает отсортированный массив строк.
func Sort(arr []string) []string {
	sort.Strings(arr)
	return arr
}

// Функция для проверки наличия строки в слайсе.
func Contains(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

// Функция для добавления строки в слайс, если она еще не присутствует.
func AddIfUnique(slice []string, s string) []string {
	if !Contains(slice, s) {
		slice = append(slice, s)
	}
	return slice
}
func RemoveString(arr []string, remove string) []string {
	// Создаем новый срез, в который будем добавлять строки
	var result []string

	// Перебираем исходный срез
	for _, item := range arr {
		// Если текущая строка не совпадает со строкой, которую нужно удалить
		if item != remove {
			// Добавляем строку в новый срез
			result = append(result, item)
		}
	}

	// Возвращаем новый срез без удаленной строки
	return result
}
func RemoveDuplicates(mainArray, excludeArray []string) []string {
	// Конвертировать excludeArray в словарь для быстрого поиска
	excludeDict := make(map[string]bool)
	for _, val := range excludeArray {
		excludeDict[val] = true
	}

	// Создать новый массив, который содержит только элементы, которые не присутствуют в excludeArray
	var resultArray []string
	for _, val := range mainArray {
		if !excludeDict[val] {
			resultArray = append(resultArray, val)
		}
	}

	return resultArray
}
func AddPrefix(strings []string, prefix ...string) []string {
	result := make([]string, len(strings))
	for i, str := range strings {
		newStr := str
		if len(prefix) > 0 {
			newStr = prefix[0] + newStr
		}
		result[i] = newStr
	}
	return result
}
