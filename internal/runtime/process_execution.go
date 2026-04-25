package runtime

import (
	"bytes"
	"os"
	"os/exec"
	"sort"
)

func RunCommand(
	directory string,
	environment map[string]string,
	name string,
	args ...string,
) (output string, err error) {
	commandPath, err := lookCommandPath(name, environment)
	if err != nil {
		return "", err
	}

	command := exec.Command(commandPath, args...)
	command.Dir = directory
	command.Env = append(os.Environ(), environmentEntries(environment)...)

	var buffer bytes.Buffer
	command.Stdout = &buffer
	command.Stderr = &buffer

	err = command.Run()
	return buffer.String(), err
}

func environmentEntries(environment map[string]string) (values []string) {
	if len(environment) == 0 {
		return nil
	}

	keys := make([]string, 0, len(environment))
	for key := range environment {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	values = make([]string, 0, len(keys))
	for _, key := range keys {
		values = append(values, key+"="+environment[key])
	}

	return values
}
