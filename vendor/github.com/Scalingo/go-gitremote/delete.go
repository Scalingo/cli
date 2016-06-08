package gitremote

import "strings"

func (config *Config) Delete(r *Remote) error {
	content, err := config.Content()
	if err != nil {
		return err
	}

	newContent := deleteFromContent(content, r)
	return config.Write(newContent)
}

func deleteFromContent(content string, r *Remote) string {
	newContent := ""
	inRemote := false
	for _, line := range strings.Split(content, "\n") {
		matches := remoteTitleRe.FindStringSubmatch(line)
		if len(matches) == 2 && r.Name == matches[1] {
			inRemote = true
		} else if inRemote && strings.HasPrefix(line, "[") {
			inRemote = false
		}

		if !inRemote {
			newContent += line + "\n"
		}
	}
	return strings.TrimRight(newContent, "\n") + "\n"
}
