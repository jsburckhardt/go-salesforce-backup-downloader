package cmd

import "strings"

func folderValidator(folderInput string) string {
	var folderOutput string

	if strings.HasSuffix(folderInput, "/") {
		folderInput = folderInput[:len(folderInput)-1]
	}

	folderOutput = folderInput
	return folderOutput
}
