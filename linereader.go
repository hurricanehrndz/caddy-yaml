package caddyyaml

import (
	"bufio"
	"bytes"
	"io"
	"regexp"
	"strings"
	"unicode"
)

var variableLineRegexp = regexp.MustCompile(`^x\-([a-zA-Z0-9\.\_\-]+)(\s*)\:`)

// commentLine checks if a line is a comment (starts with # after trimming whitespace).
func commentLine(line string) bool {
	return strings.HasPrefix(strings.TrimSpace(line), "#")
}

// topLevelLine checks if a line is a top-level YAML key (not indented, not a comment, not an x- variable).
func topLevelLine(line string) bool {
	trimmed := strings.TrimLeftFunc(line, unicode.IsSpace)
	return trimmed != "" && trimmed == line && !commentLine(line) && !strings.HasPrefix(line, "x-")
}

// variableLine checks if a line declares an x- prefixed variable and returns the variable name.
func variableLine(line string) (variable string, found bool) {
	matched := variableLineRegexp.FindStringSubmatch(line)
	if len(matched) > 0 {
		return matched[1], true
	}
	return "", false
}

// extractVariables extracts variables from the body.
func extractVariables(body []byte) ([]byte, error) {
	var variablesBuffer bytes.Buffer
	reader := bufio.NewReader(bytes.NewReader(body))
	inVariable := false

	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return nil, err
		}

		// Check if this is a variable declaration line
		if _, ok := variableLine(line); ok {
			inVariable = true
			variablesBuffer.WriteString(line)
		} else if inVariable {
			// If we're in a variable, check if this line belongs to it
			if topLevelLine(line) {
				// Hit a new top-level section, stop collecting
				inVariable = false
			} else {
				// Still part of the variable (indented/empty/comment)
				variablesBuffer.WriteString(line)
			}
		}

		if err == io.EOF {
			break
		}
	}

	return variablesBuffer.Bytes(), nil
}
