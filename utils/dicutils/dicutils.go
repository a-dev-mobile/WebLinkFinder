package dicutils

/// Function to check if the value is false in the dictionary
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
func GetKeysWithTrue(dict map[string]bool) []string {
	var keys []string

	for key, value := range dict {
		if value {
			keys = append(keys, key)
		}
	}

	return keys

}
func AddToMapIfNotExist(myMap map[string]bool, key string, value bool) {
	// Check if key already exists in map
	if _, exists := myMap[key]; !exists {
		myMap[key] = value
	}
}
