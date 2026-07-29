package services

import (
	"regexp"
	"strings"
)

var markdownLinkPattern = regexp.MustCompile(`\[([^\]]+)\]\((https?://[^)\s]+)\)`)

// normalizeAssistantAnswer turns a model response into portable plain text that
// renders consistently in SwiftUI, Jetpack Compose and the React client.
func normalizeAssistantAnswer(raw string) string {
	raw = strings.ReplaceAll(raw, "\r\n", "\n")
	raw = strings.ReplaceAll(raw, "\r", "\n")
	sourceLines := strings.Split(raw, "\n")
	lines := make([]string, 0, len(sourceLines))

	for index := 0; index < len(sourceLines); {
		line := strings.TrimSpace(sourceLines[index])
		if isMarkdownTableRow(line) && index+1 < len(sourceLines) &&
			isMarkdownTableSeparator(sourceLines[index+1]) {
			headers := markdownTableCells(line)
			index += 2
			for index < len(sourceLines) && isMarkdownTableRow(strings.TrimSpace(sourceLines[index])) {
				cells := markdownTableCells(sourceLines[index])
				if !isMarkdownTableSeparator(sourceLines[index]) {
					lines = append(lines, formatTableRecord(headers, cells))
				}
				index++
			}
			continue
		}

		lines = append(lines, normalizePortableLine(line))
		index++
	}

	lines = addListSpacing(lines)
	return strings.TrimSpace(strings.Join(collapseBlankLines(lines), "\n"))
}

func normalizePortableLine(line string) string {
	line = strings.TrimSpace(line)
	if line == "" {
		return ""
	}
	if line == "```" || strings.HasPrefix(line, "```") {
		return ""
	}
	if strings.Trim(line, "-*_ ") == "" {
		return ""
	}

	for strings.HasPrefix(line, "#") {
		line = strings.TrimSpace(strings.TrimPrefix(line, "#"))
	}
	if strings.HasPrefix(line, "> ") {
		line = strings.TrimSpace(strings.TrimPrefix(line, "> "))
	}
	for _, prefix := range []string{"- ", "* ", "+ "} {
		if strings.HasPrefix(line, prefix) {
			line = "• " + strings.TrimSpace(strings.TrimPrefix(line, prefix))
			break
		}
	}

	line = strings.ReplaceAll(line, "**", "")
	line = strings.ReplaceAll(line, "__", "")
	line = markdownLinkPattern.ReplaceAllString(line, "$1 ($2)")
	return strings.TrimSpace(line)
}

func isMarkdownTableRow(line string) bool {
	line = strings.TrimSpace(line)
	return strings.HasPrefix(line, "|") && strings.HasSuffix(line, "|") &&
		strings.Count(line, "|") >= 3
}

func isMarkdownTableSeparator(line string) bool {
	if !isMarkdownTableRow(strings.TrimSpace(line)) {
		return false
	}
	cells := markdownTableCells(line)
	if len(cells) == 0 {
		return false
	}
	for _, cell := range cells {
		value := strings.TrimSpace(strings.ReplaceAll(cell, ":", ""))
		if len(value) < 3 || strings.Trim(value, "-") != "" {
			return false
		}
	}
	return true
}

func markdownTableCells(line string) []string {
	line = strings.Trim(strings.TrimSpace(line), "|")
	rawCells := strings.Split(line, "|")
	cells := make([]string, 0, len(rawCells))
	for _, cell := range rawCells {
		cell = strings.TrimSpace(cell)
		cell = strings.ReplaceAll(cell, "**", "")
		cell = strings.ReplaceAll(cell, "__", "")
		cell = markdownLinkPattern.ReplaceAllString(cell, "$1 ($2)")
		cells = append(cells, cell)
	}
	return cells
}

func formatTableRecord(headers, cells []string) string {
	if len(cells) == 0 {
		return ""
	}
	if len(cells) == 1 {
		return "• " + cells[0]
	}

	var details []string
	for index := 1; index < len(cells); index++ {
		if index < len(headers) && headers[index] != "" {
			details = append(details, headers[index]+": "+cells[index])
		} else {
			details = append(details, cells[index])
		}
	}
	return "• " + cells[0] + ": " + strings.Join(details, "; ")
}

func addListSpacing(lines []string) []string {
	result := make([]string, 0, len(lines)+4)
	for index, line := range lines {
		isBullet := strings.HasPrefix(line, "• ")
		previousBullet := len(result) > 0 && strings.HasPrefix(result[len(result)-1], "• ")
		if isBullet && len(result) > 0 && result[len(result)-1] != "" && !previousBullet {
			result = append(result, "")
		}
		if !isBullet && line != "" && previousBullet {
			result = append(result, "")
		}
		result = append(result, line)

		if index == len(lines)-1 && isBullet {
			result = append(result, "")
		}
	}
	return result
}

func collapseBlankLines(lines []string) []string {
	result := make([]string, 0, len(lines))
	previousBlank := true
	for _, line := range lines {
		blank := strings.TrimSpace(line) == ""
		if blank && previousBlank {
			continue
		}
		result = append(result, line)
		previousBlank = blank
	}
	for len(result) > 0 && strings.TrimSpace(result[len(result)-1]) == "" {
		result = result[:len(result)-1]
	}
	return result
}
