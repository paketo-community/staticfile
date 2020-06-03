package staticfile

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type BuildpackYMLParser struct{}

func NewBuildpackYMLParser() BuildpackYMLParser {
	return BuildpackYMLParser{}
}

func (p BuildpackYMLParser) Parse(path string) (Config, error) {
	return parse(path)
}

func (p BuildpackYMLParser) ValidConfig(path string) (bool, error) {
	config, err := parse(path)
	if err != nil {
		return false, err
	}

	if config.Nginx != nil {
		return true, nil
	}

	return false, nil

}

func parse(path string) (Config, error) {
	var buildpack struct {
		Staticfile Config `yaml:"staticfile"`
	}

	file, err := os.Open(path)
	if os.IsNotExist(err) {
		return Config{}, nil
	} else if err != nil {
		return Config{}, fmt.Errorf("unable to open buildpack.yml: %q", err)
	}

	err = yaml.NewDecoder(file).Decode(&buildpack)
	if err != nil {
		return Config{}, fmt.Errorf("unable to parse buildpack.yml: %q", err)
	}

	return buildpack.Staticfile, err

}
