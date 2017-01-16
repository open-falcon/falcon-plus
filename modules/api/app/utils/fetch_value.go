package utils

func GetValue(maplist []map[string]interface{}, key string) (result []interface{}) {
	for _, v := range maplist {
		result = append(result, v[key])
	}
	return
}
