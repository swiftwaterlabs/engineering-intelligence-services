package core

import "encoding/json"

func MapFromJson(toMap string, target interface{}) error {
	err := json.Unmarshal([]byte(toMap), target)

	return err
}
