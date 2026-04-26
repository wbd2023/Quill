package runtime

import (
	"archive/tar"
	"fmt"
	"path"
	"strings"
)

func validateShellcheckArchiveEntry(
	header *tar.Header,
	version string,
) (name string, err error) {
	switch header.Typeflag {
	case tar.TypeSymlink, tar.TypeLink:
		return "", fmt.Errorf("shellcheck archive contains link entry %q", header.Name)
	}

	rawName := header.Name
	if header.Typeflag == tar.TypeDir {
		rawName = strings.TrimSuffix(rawName, "/")
	}

	name = path.Clean(rawName)
	if name == "." ||
		name != rawName ||
		path.IsAbs(rawName) ||
		strings.HasPrefix(name, "../") ||
		strings.Contains(name, "/../") {
		return "", fmt.Errorf("unsafe shellcheck archive path %q", header.Name)
	}

	root := shellcheckArchiveRoot(version)
	switch name {
	case root,
		path.Join(root, "LICENSE.txt"),
		path.Join(root, "README.txt"),
		path.Join(root, "shellcheck"):
		return name, nil
	default:
		return "", fmt.Errorf("unexpected shellcheck archive entry %q", header.Name)
	}
}

func shellcheckArchiveRoot(version string) (root string) {
	return "shellcheck-v" + version
}
