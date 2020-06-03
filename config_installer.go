package staticfile

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/paketo-buildpacks/packit"
)

type ConfigInstaller struct{}

func NewConfigInstaller() ConfigInstaller {
	return ConfigInstaller{}
}

var ConfigMap map[string]string = map[string]string{
	"nginx": "nginx.conf",
}

func (ci ConfigInstaller) Execute(context packit.BuildContext, templConfig Config) error {
	conf, ok := ConfigMap["nginx"]
	if !ok {
		return fmt.Errorf("%q does not have a corresponding config file: supported servers are: [%s]", "nginx", strings.Join(supportedServers(), ","))
	}

	confSource := filepath.Join(context.CNBPath, "server_configs", conf)
	confDest := filepath.Join(context.WorkingDir, conf)

	confContents, err := generateConf(confSource, templConfig)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(confDest, []byte(confContents), 0644)

	if err != nil {
		return fmt.Errorf("could not write nginx conf file to working dir: %w", err)
	}

	return nil

}

func supportedServers() []string {
	servers := make([]string, len(ConfigMap))

	i := 0
	for k := range ConfigMap {
		servers[i] = k
		i++
	}

	return servers
}

func generateConf(configPath string, templConfig Config) (string, error) {
	buffer := new(bytes.Buffer)

	config, err := ioutil.ReadFile(configPath)
	if err != nil {
		return "", fmt.Errorf("could not find template file: %w", err)
	}

	t := template.Must(template.New("nginx.conf").Delims("$((", "))").Parse(string(config)))

	err = t.Execute(buffer, templConfig.Nginx)
	if err != nil {
		// not tested
		return "", fmt.Errorf("templating failed: %w", err)
	}
	return buffer.String(), nil
}
