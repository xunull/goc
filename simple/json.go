package simple

import "encoding/json"

func JsonMarshal(i interface{}) ([]byte, error) {
	return json.Marshal(i)
}

func JsonUnmarshal(data []byte, i interface{}) error {
	return json.Unmarshal(data, i)
}
