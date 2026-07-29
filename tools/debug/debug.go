package debug

import (
	"encoding/json"
	"fmt"
)

// DumpJson DebugDumpJson 打印json
func DumpJson(v map[string]any, pretty bool) {
	if pretty {
		jsonPretty, err := json.MarshalIndent(v, "", "  ")
		if err != nil {
			fmt.Println("json MarshalIndent error:", err)
			return
		}
		fmt.Println(string(jsonPretty))
		return
	}
	jsonData, err := json.Marshal(v)
	if err != nil {
		fmt.Println("json Marshal error:", err)
		return
	}
	fmt.Println(string(jsonData))
	return
}
