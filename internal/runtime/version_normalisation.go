package runtime

import (
	"fmt"
	"strconv"
	"strings"
)

func parseShellcheckVersion(output string) (version string, err error) {
	for line := range strings.SplitSeq(output, "\n") {
		line = strings.TrimSpace(line)
		if after, ok := strings.CutPrefix(line, "version:"); ok {
			return strings.TrimSpace(after), nil
		}
	}

	return "", fmt.Errorf("could not parse shellcheck version")
}

func parseSingleTokenVersion(output string) (version string, err error) {
	fields := strings.Fields(output)
	if len(fields) == 0 {
		return "", fmt.Errorf("could not parse version output")
	}

	return strings.TrimPrefix(fields[0], "v"), nil
}

func matchesPinnedVersion(actual string, pinned string) (matches bool) {
	return normaliseVersion(actual) == normaliseVersion(pinned)
}

func normaliseVersion(version string) (normalised string) {
	version = strings.TrimPrefix(strings.TrimSpace(version), "v")
	end := len(version)
	for index, rune := range version {
		if (rune < '0' || rune > '9') && rune != '.' {
			end = index
			break
		}
	}

	version = version[:end]
	if version == "" {
		return ""
	}

	parts := strings.Split(version, ".")
	for _, part := range parts {
		if _, err := strconv.Atoi(part); err != nil {
			return ""
		}
	}

	return strings.Join(parts, ".")
}
