package main

import (
	"os"

	"gopkg.in/yaml.v2"
)

type BuildpackYMLParser struct{}

func NewBuildpackYMLParser() BuildpackYMLParser {
	return BuildpackYMLParser{}
}

func (p BuildpackYMLParser) Parse(path string) (string, error) {
	var buildpack struct {
		Staticfile struct {
			Server string `yaml:"server"`
		} `yaml:"staticfile"`
	}
	file, err := os.Open(path)
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}
	defer file.Close()

	if !os.IsNotExist(err) {
		err = yaml.NewDecoder(file).Decode(&buildpack)
		if err != nil {
			panic(err)
		}
	}

	return buildpack.Staticfile.Server, err
}
