package utils

import "encoding/json"

type Object map[string]any

// its the caller responsability to not call toJsonString on an object that may not be marshaled
func ToJsonString(obj any) string {
	bs, _ := json.MarshalIndent(obj, "\t", " ")
	return string(bs)
}
