package installer

import "fmt"

// archiveHashes holds the recorded SHA-256 hashes for archive-installed tools,
// keyed by tool ID then GOOS/GOARCH. Phase 3 replaces this with quill.lock.
var archiveHashes = map[string]map[string]string{
	"shellcheck": {
		"darwin/amd64": "ef27684f23279d112d8ad84e0823642e43f838993bbb8c0963db9b58a90464c2",
		"darwin/arm64": "bbd2f14826328eee7679da7221f2bc3afb011f6a928b848c80c321f6046ddf81",
		"linux/amd64":  "6c881ab0698e4e6ea235245f22832860544f17ba386442fe7e9d629f8cbedf87",
		"linux/arm64":  "324a7e89de8fa2aed0d0c28f3dab59cf84c6d74264022c00c22af665ed1a09bb",
	},
}

func archiveHashFor(toolID string, goos string, goarch string) (hash string, err error) {
	byPlatform, ok := archiveHashes[toolID]
	if !ok {
		return "", fmt.Errorf("no recorded hashes for tool %s", toolID)
	}

	hash, ok = byPlatform[goos+"/"+goarch]
	if !ok {
		return "", fmt.Errorf("no recorded hash for %s on %s/%s", toolID, goos, goarch)
	}

	return hash, nil
}
