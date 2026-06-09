package render

import (
	"bytes"
	"os/exec"
	"strings"
)

func Something(markdown string) string {
	cmd := exec.Command("bat", "--language", "md", "--style=plain", "--color=always", "--paging=never")
	cmd.Stdin = strings.NewReader(markdown)

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		print(err)
		return markdown
	}

	return stdout.String()
}

