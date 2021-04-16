package simple

import "gopkg.in/yaml.v2"

func YamlUnmarshal(bs []byte, i interface{}) error {
	return yaml.Unmarshal(bs, i)
}

func Yaml(in interface{}) ([]byte, error) {
	return yaml.Marshal(in)
}
