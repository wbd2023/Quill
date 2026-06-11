package text

import "ciphera/tools/internal/policy"

func fileSets() (fileSets policy.FileSets) {
	fileSets = append(fileSets, policy.FileSetConfig{
		Name: "line_length",
		Exclude: policy.FileSetExclude{
			Files: []string{"go.sum", "package-lock.json"},
		},
	})
	fileSets = append(fileSets, policy.FileSetConfig{
		Name: "spelling",
		Exclude: policy.FileSetExclude{
			Extensions: []string{".go"},
			Files: []string{
				"COPYING",
				"COPYING.*",
				"LICENSE",
				"LICENSE.*",
				"NOTICE",
				"NOTICE.*",
				"package-lock.json",
			},
		},
	})
	return fileSets
}
