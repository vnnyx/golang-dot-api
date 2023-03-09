package util

import "encoding/json"

func ToMap(s interface{}) (m map[string]interface{}, err error) {
	b, err := json.Marshal(s)
	if err != nil {
		return m, err
	}
	err = json.Unmarshal(b, &m)
	if err != nil {
		return m, err
	}
	return m, nil
}
