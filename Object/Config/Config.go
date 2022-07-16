package Config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

func ReadConfig(path string) Yaml {
	var conf Yaml
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return conf
	}
	if err = yaml.Unmarshal(yamlFile, &conf); err != nil {
		fmt.Println(err)
		return conf
	}
	return conf
}
