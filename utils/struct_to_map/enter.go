package struct_to_map

import (
	"Blog_server/models/enum"
	"github.com/fatih/structs"
	"github.com/sirupsen/logrus"
)

func StructToMap(obj interface{}) map[string]interface{} {
	m := structs.Map(obj)
	data := DeleteEmpty(m)
	logrus.Infof("data: %v", data)
	return data
}

func DeleteEmpty(m map[string]interface{}) map[string]interface{} {
	var data = make(map[string]interface{}, 0)
	for key, v := range m {
		switch val := v.(type) {
		case string:
			if val != "" {
				data[key] = val
			}
		case int:
			if val != 0 {
				data[key] = val
			}
		case uint:
			if val != 0 {
				data[key] = val
			}
		case int64:
			if val != 0 {
				data[key] = val
			}
		case []string:
			if val != nil {
				data[key] = val
			}
		case enum.Array:
			if len(val) != 0 {
				data[key] = val
			}
		default:
			data[key] = v
		}

	}
	return data
}
