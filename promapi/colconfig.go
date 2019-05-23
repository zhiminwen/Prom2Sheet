package promapi

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

type Column struct {
	Name  string `yaml:"name"`
	Type  string `yaml:"type"`
	Value string `yaml:"value"`
}

type Sheet struct {
	Name    string   `yaml:"name"`
	Query   string   `yaml:"query"`
	Columns []Column `yaml:"columns"`
}
type XLSConfig struct {
	Sheets []Sheet `yaml:"sheets,omitempty"`
}

func ParseSheetYaml(yamlfile string) (*XLSConfig, error) {
	content, err := ioutil.ReadFile(yamlfile)
	if err != nil {
		log.Printf("Failed to read yaml file:%v", err)
		return nil, err
	}

	var config XLSConfig
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		log.Printf("Failed to parse yaml content:%v", err)
		return nil, err
	}

	return &config, nil
}
