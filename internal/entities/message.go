package entities

import "strings"

type Message struct {
	Text      string
	ParseMode ParseMode
	Token     string
	ChatIds   []int64
}

type ParseMode string

var (
	parseModeMap = map[string]ParseMode{
		"markdownv2": MarkdownV2,
		"markdown":   Markdown,
		"html":       HTML,
	}
)

const (
	Undefined  ParseMode = "Undefined"
	MarkdownV2 ParseMode = "MarkdownV2"
	Markdown   ParseMode = "Markdown"
	HTML       ParseMode = "HTML"
)

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

func ParseString(str string) ParseMode {
	c, ok := parseModeMap[strings.ToLower(str)]
	if !ok {
		return Undefined
	}

	return c
}
