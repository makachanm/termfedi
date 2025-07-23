package layer

import (
	"encoding/json"
	"fmt"
)

func unmarshallJSON[T any](obj *T, d []byte) {
	err := json.Unmarshal(d, obj)
	if err != nil {
		fmt.Println("BODY: ", string(d))
		panic(err)
	}
}

func marshallJSON[T any](obj *T) []byte {
	d, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}

	return d
}
