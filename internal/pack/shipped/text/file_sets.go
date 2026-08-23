package text

import "github.com/wbd2023/quill/internal/profile"

func fileSets() (fileSets profile.FileSets) {
	fileSets = append(fileSets, profile.FileSetConfig{
		Name: "line_length",
		Exclude: profile.FileSetExclude{
			Files: []string{"go.sum", "package-lock.json"},
		},
	})
	fileSets = append(fileSets, profile.FileSetConfig{
		Name: "spelling",
		Exclude: profile.FileSetExclude{
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
