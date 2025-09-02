package layer

import (
	"encoding/json"
)

func unmarshallJSON[T any](obj *T, d []byte) bool {
	err := json.Unmarshal(d, obj)
	if err != nil {
		return false
	}

	return true
}

func marshallJSON[T any](obj *T) []byte {
	d, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}

	return d
}
