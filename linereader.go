package caddyyaml

import (
	"bufio"
	"bytes"
	"io"
	"regexp"
	"strings"
	"unicode"
)

var extensionLineRegexp = regexp.MustCompile(`^x\-([a-zA-Z0-9\.\_\-]+)(\s*)\:`)

// commentLine checks if a line is a comment (starts with # after trimming whitespace).
func commentLine(line string) bool {
	return strings.HasPrefix(strings.TrimSpace(line), "#")
}

// isTopLevelKey checks if a line is a top-level YAML key (not indented, not a comment, not an x- extension field).
func isTopLevelKey(line string) bool {
	trimmed := strings.TrimLeftFunc(line, unicode.IsSpace)
	hasNoIndentation := trimmed != "" && trimmed == line

	isComment := commentLine(line)
	isExtensionField := strings.HasPrefix(line, "x-")

	return hasNoIndentation && !isComment && !isExtensionField
}

// isExtensionLine checks if a line declares an x- prefixed extension field and returns the field name.
func isExtensionLine(line string) (fieldName string, found bool) {
	matched := extensionLineRegexp.FindStringSubmatch(line)
	if len(matched) > 0 {
		return matched[1], true
	}
	return "", false
}

// extractRawExtensions extracts x- prefixed extension fields from the body.
func extractRawExtensions(body []byte) ([]byte, error) {
	var extensionsBuffer bytes.Buffer
	reader := bufio.NewReader(bytes.NewReader(body))
	inExtension := false

	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return nil, err
		}

		inExtension = collectExtensionLine(line, &extensionsBuffer, inExtension)

		if err == io.EOF {
			break
		}
	}

	return extensionsBuffer.Bytes(), nil
}

// collectExtensionLine collects a line into the buffer if it's part of an extension field definition.
// Returns the updated inExtension state.
func collectExtensionLine(line string, buffer *bytes.Buffer, inExtension bool) bool {
	// Check if this is an extension field declaration line
	if _, ok := isExtensionLine(line); ok {
		buffer.WriteString(line)
		return true
	}

	// If we're in an extension field, check if this line belongs to it
	if inExtension {
		if isTopLevelKey(line) {
			// Hit a new top-level section, stop collecting
			return false
		}
		// Still part of the extension field (indented/empty/comment)
		buffer.WriteString(line)
	}

	return inExtension
}
