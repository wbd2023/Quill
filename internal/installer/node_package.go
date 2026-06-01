package installer

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func validatePackageMetadata(packagePath string, lockPath string) (err error) {
	packageName, err := readPackageName(packagePath)
	if err != nil {
		return err
	}

	lockName, err := readPackageName(lockPath)
	if err != nil {
		return err
	}

	if lockName != packageName {
		return fmt.Errorf(
			"package-lock name %q does not match package.json name %q",
			lockName,
			packageName,
		)
	}

	return nil
}

func readPackageName(path string) (name string, err error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	var document struct {
		Name string `json:"name"`
	}
	if err = json.Unmarshal(contents, &document); err != nil {
		return "", err
	}

	if document.Name == "" {
		return "", fmt.Errorf("%s does not define a package name", filepath.Base(path))
	}

	return document.Name, nil
}
