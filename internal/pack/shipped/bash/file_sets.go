package bash

import "github.com/wbd2023/quill/internal/policy"

func fileSets() (fileSets policy.FileSets) {
	return append(fileSets, policy.FileSetConfig{
		Name: "bash",
		Include: policy.FileSetInclude{
			Extensions: []string{".sh"},
		},
	})
}
