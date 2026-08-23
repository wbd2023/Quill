package bash

import "github.com/wbd2023/quill/internal/profile"

func fileSets() (fileSets profile.FileSets) {
	return append(fileSets, profile.FileSetConfig{
		Name: "bash",
		Include: profile.FileSetInclude{
			Extensions: []string{".sh"},
		},
	})
}
