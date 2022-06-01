package staticfile

import (
	"fmt"
	"os"
	"path/filepath"
)

type ProfileDWriter struct{}

func NewProfileDWriter() ProfileDWriter {
	return ProfileDWriter{}
}

func (p ProfileDWriter) WriteStartLoggingScript(profileDPath string) error {
	return writeScript(profileDPath, "01_start_logging.sh", StartLoggingContents)
}

func (p ProfileDWriter) WriteInitScript(profileDPath string) error {
	return writeScript(profileDPath, "00_staticfile.sh", InitScriptContents)
}

func writeScript(dir, file, contents string) error {
	err := os.MkdirAll(dir, 0744)
	if err != nil {
		return fmt.Errorf("could not create the profile.d layer directory: %v", err)
	}

	scriptPath := filepath.Join(dir, file)
	err = os.WriteFile(scriptPath, []byte(contents), 0744)

	if err != nil {
		return fmt.Errorf("could not write the %s script: %v", file, err)
	}

	return nil
}
