package caddyyaml

import (
	"bufio"
	"bytes"
	"io"
	"regexp"
	"strings"
	"unicode"
)

// commentLine checks if a line is a comment (starts with # after trimming whitespace).
func commentLine(line string) bool {
	return strings.HasPrefix(strings.TrimSpace(line), "#")
}

// isTopLevelKey checks if a line is a top-level YAML key (not indented, not a comment).
func isTopLevelKey(line string) bool {
	trimmed := strings.TrimLeftFunc(line, unicode.IsSpace)
	hasNoIndentation := trimmed != "" && trimmed == line

	isComment := commentLine(line)

	return hasNoIndentation && !isComment
}

// extractAllMatchingTopLevelSections extracts all top-level YAML sections matching the given pattern from the body.
// It returns all extracted sections concatenated and the remaining body with sections removed.
// The pattern should be a compiled regexp that matches the top-level key lines.
func extractAllMatchingTopLevelSections(body []byte, pattern *regexp.Regexp) (sections []byte, remaining []byte) {
	var sectionsBuffer bytes.Buffer
	var remainingBuffer bytes.Buffer
	reader := bufio.NewReader(bytes.NewReader(body))
	inSection := false

	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return body, nil
		}

		inSection = processLine(line, pattern, &sectionsBuffer, &remainingBuffer, inSection)

		if err == io.EOF {
			break
		}
	}

	return sectionsBuffer.Bytes(), remainingBuffer.Bytes()
}

// processLine processes a single line for extractAllMatchingTopLevelSections.
// Returns the updated inSection state.
func processLine(line string, pattern *regexp.Regexp, sectionsBuffer, remainingBuffer *bytes.Buffer, inSection bool) bool {
	if pattern.MatchString(line) {
		sectionsBuffer.WriteString(line)
		return true
	}

	if inSection {
		return handleSectionLine(line, sectionsBuffer, remainingBuffer)
	}

	remainingBuffer.WriteString(line)
	return false
}

// handleSectionLine handles a line when already inside a section for extractAllMatchingTopLevelSections.
// Returns whether we're still in the section.
func handleSectionLine(line string, sectionsBuffer, remainingBuffer *bytes.Buffer) bool {
	if isTopLevelKey(line) {
		remainingBuffer.WriteString(line)
		return false
	}
	sectionsBuffer.WriteString(line)
	return true
}
