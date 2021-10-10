package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

type Config struct {
	Port string `yaml:"port"`
	Mysql struct {
		Host         string `yaml:"host"`
		Port         string `yaml:"port"`
		Username     string `yaml:"username"`
		Password     string `yaml:"password"`
		DatabaseName string `yaml:"dbname"`
	}
}

var Configuration Config

func init() {
	configBytes, err := ioutil.ReadFile("pkg/config/conf.yaml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(configBytes, &Configuration)
	if err != nil {
		panic(err)
	}
	log.Printf("Successfully parsed config: %+v\n", Configuration)
}
