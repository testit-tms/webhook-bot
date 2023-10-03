package entities

import "strings"

// Message represents a message to be sent to one or more chats.
type Message struct {
	Text      string
	ParseMode ParseMode
	Token     string
	ChatIds   []int64
}

// ParseMode represents the parsing mode for a message.
type ParseMode string

var (
	parseModeMap = map[string]ParseMode{
		"markdownv2": MarkdownV2,
		"markdown":   Markdown,
		"html":       HTML,
	}
)

const (
	// Undefined represents an undefined parsing mode.
	Undefined  ParseMode = "Undefined"
	// MarkdownV2 represents the MarkdownV2 parsing mode.
	MarkdownV2 ParseMode = "MarkdownV2"
	// Markdown represents the Markdown parsing mode.
	Markdown   ParseMode = "Markdown"
	// HTML represents the HTML parsing mode.
	HTML       ParseMode = "HTML"
)

// String returns the string representation of the ParseMode.
func (pm ParseMode) String() string {
	switch pm {
	case MarkdownV2:
		return "MarkdownV2"
	case Markdown:
		return "Markdown"
	case HTML:
		return "HTML"
	default:
		return ""
	}
}

// ParseString returns the ParseMode for the given string.
func ParseString(str string) ParseMode {
	c, ok := parseModeMap[strings.ToLower(str)]
	if !ok {
		return Undefined
	}

	return c
}
