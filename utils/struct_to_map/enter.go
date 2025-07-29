package struct_to_map

import (
	"fmt"
	"github.com/fatih/structs"
)

func StructToMap(obj interface{}) map[string]interface{} {
	m := structs.Map(obj)
	for k, _ := range m {
		if m[k] == "" {
			delete(m, k)
		}
	}
	fmt.Println(m)
	return m

}
