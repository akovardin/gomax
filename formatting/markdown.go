package formatting

import (
	"regexp"

	"github.com/akovardin/gomax/types"
)

type Formatter struct{}

func NewFormatter() *Formatter {
	return &Formatter{}
}

func utf16Len(s string) int {
	count := 0
	for _, r := range s {
		if r >= 0x10000 {
			count += 2
		} else {
			count += 1
		}
	}
	return count
}

func (f *Formatter) FormatMarkdown(text string) (string, []types.Element) {
	var elements []types.Element

	type pattern struct {
		re     *regexp.Regexp
		elType string
	}

	patterns := []pattern{
		{regexp.MustCompile(`\*\*(.+?)\*\*`), "STRONG"},
		{regexp.MustCompile(`__(.+?)__`), "UNDERLINE"},
		{regexp.MustCompile(`~~(.+?)~~`), "STRIKETHROUGH"},
		{regexp.MustCompile("`([^`]+)`"), "MONOSPACED"},
		{regexp.MustCompile(`\b_(.+?)_\b`), "EMPHASIZED"},
		{regexp.MustCompile(`\*(.+?)\*`), "EMPHASIZED"},
		{regexp.MustCompile(`\[(.+?)\]\((.+?)\)`), "LINK"},
	}

	result := text
	for _, p := range patterns {
		for {
			loc := p.re.FindStringSubmatchIndex(result)
			if loc == nil {
				break
			}

			fullStart := loc[0]
			fullEnd := loc[1]
			contentStart := loc[2]
			contentEnd := loc[3]

			content := result[contentStart:contentEnd]
			length := utf16Len(content)
			beforeUtf16 := utf16Len(result[:contentStart])

			from := beforeUtf16
			l := length

			el := types.Element{
				Type:   p.elType,
				From:   &from,
				Length: &l,
			}

			if p.elType == "LINK" && len(loc) >= 6 {
				urlStart := loc[4]
				urlEnd := loc[5]
				url := result[urlStart:urlEnd]
				el.Attributes = &types.ElementAttributes{URL: &url}
			}

			elements = append(elements, el)
			result = result[:fullStart] + content + result[fullEnd:]
		}
	}

	codeBlockRe := regexp.MustCompile("(?s)```(.+?)```")
	for {
		loc := codeBlockRe.FindStringSubmatchIndex(result)
		if loc == nil {
			break
		}
		contentStart := loc[2]
		contentEnd := loc[3]
		content := result[contentStart:contentEnd]
		length := utf16Len(content)
		beforeUtf16 := utf16Len(result[:contentStart])

		from := beforeUtf16
		l := length

		el := types.Element{
			Type:   "CODE",
			From:   &from,
			Length: &l,
		}
		elements = append(elements, el)

		result = result[:loc[0]] + content + result[loc[1]:]
	}

	headingRe := regexp.MustCompile(`(?m)^#{1,3}\s+(.+)$`)
	for {
		loc := headingRe.FindStringSubmatchIndex(result)
		if loc == nil {
			break
		}
		contentStart := loc[2]
		contentEnd := loc[3]
		content := result[contentStart:contentEnd]
		length := utf16Len(content)
		beforeUtf16 := utf16Len(result[:contentStart])

		from := beforeUtf16
		l := length

		el := types.Element{
			Type:   "HEADING",
			From:   &from,
			Length: &l,
		}
		elements = append(elements, el)

		result = result[:loc[0]] + content + result[loc[1]:]
	}

	quoteRe := regexp.MustCompile(`(?m)^>\s+(.+)$`)
	for {
		loc := quoteRe.FindStringSubmatchIndex(result)
		if loc == nil {
			break
		}
		contentStart := loc[2]
		contentEnd := loc[3]
		content := result[contentStart:contentEnd]
		length := utf16Len(content)
		beforeUtf16 := utf16Len(result[:contentStart])

		from := beforeUtf16
		l := length

		el := types.Element{
			Type:   "QUOTE",
			From:   &from,
			Length: &l,
		}
		elements = append(elements, el)

		result = result[:loc[0]] + content + result[loc[1]:]
	}

	return result, elements
}
