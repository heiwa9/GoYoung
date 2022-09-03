package lib

import (
	"encoding/json"
	"os"
)

func WriteJsonFile[T any](filepath string, model *T) (err error) {
	data, err := json.Marshal(model)
	if err != nil {
		return
	}
	err = os.WriteFile(filepath, data, 0664)
	return
}

func ReadJsonBind[T any](filepath string, model *T) (err error) {
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		return
	}
	err = json.Unmarshal(bytes, model)
	return
}
