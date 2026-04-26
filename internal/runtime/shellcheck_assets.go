package runtime

import "fmt"

var shellcheckAssets = map[string]shellcheckAsset{
	"darwin/amd64": {
		Name:   "darwin.x86_64",
		SHA256: "ef27684f23279d112d8ad84e0823642e43f838993bbb8c0963db9b58a90464c2",
	},
	"darwin/arm64": {
		Name:   "darwin.aarch64",
		SHA256: "bbd2f14826328eee7679da7221f2bc3afb011f6a928b848c80c321f6046ddf81",
	},
	"linux/amd64": {
		Name:   "linux.x86_64",
		SHA256: "6c881ab0698e4e6ea235245f22832860544f17ba386442fe7e9d629f8cbedf87",
	},
	"linux/arm64": {
		Name:   "linux.aarch64",
		SHA256: "324a7e89de8fa2aed0d0c28f3dab59cf84c6d74264022c00c22af665ed1a09bb",
	},
}

type shellcheckAsset struct {
	Name   string
	SHA256 string
}

func shellcheckAssetName(goos string, goarch string) (name string, err error) {
	asset, err := shellcheckAssetFor(goos, goarch)
	if err != nil {
		return "", err
	}

	return asset.Name, nil
}

func shellcheckAssetFor(goos string, goarch string) (asset shellcheckAsset, err error) {
	asset, found := shellcheckAssets[goos+"/"+goarch]
	if !found {
		return shellcheckAsset{}, fmt.Errorf("unsupported shellcheck platform: %s/%s", goos, goarch)
	}

	if asset.SHA256 == "" {
		return shellcheckAsset{}, fmt.Errorf("missing shellcheck checksum for %s/%s", goos, goarch)
	}

	return asset, nil
}
