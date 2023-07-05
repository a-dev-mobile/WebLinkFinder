package dicutils

/// Функция для проверки наличия значения false в словаре
func CheckForFalse(dict map[string]bool) bool {
	for _, value := range dict {

		return !value

	}
	return false
}
func GetKeysWithFalse(dict map[string]bool) []string {
	var keys []string

	for key, value := range dict {
		if !value {
			keys = append(keys, key)
		}
	}

	return keys

}
func AddToMapIfNotExist(myMap map[string]bool, key string, value bool) {
	// Проверить, существует ли ключ уже в map
	if _, exists := myMap[key]; !exists {
		myMap[key] = value
	}
}
