package layer

import (
	"encoding/json"
	"fmt"
)

func unmarshallJSON[T any](obj *T, d []byte) bool {
	err := json.Unmarshal(d, obj)
	if err != nil {
		//fmt.Println("BODY: ", string(d))
		fmt.Println("Error unmarshalling JSON:", err)
		return false
	}

	return true
}

func marshallJSON[T any](obj *T) []byte {
	d, err := json.Marshal(obj)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
	}

	return d
}
