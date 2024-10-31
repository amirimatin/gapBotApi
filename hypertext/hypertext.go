package hypertext

import "fmt"

func Bold(str string) string {
	return fmt.Sprintf("<b>%s</b>", str)
}

func Italic(str string) string {
	return fmt.Sprintf("<i>%s</i>\n", str)
}

func UnderLine(str string) string {
	return fmt.Sprintf("<u>%s</u>", str)
}

func Specifies(str string) string {
	return fmt.Sprintf("<s>%s</s>", str)
}

func InlineQuote(str string) string {
	return fmt.Sprintf("`%s`", str)
}

func Quote(str string) string {
	return fmt.Sprintf("```\n%s\n```", str)
}

func Colorize(str string, hexColor string) string {
	return fmt.Sprintf("<color%s>%s</color>\n", hexColor, str)
}
func Enter() string {
	return "\n"
}
func DoubleEnter() string {
	return "\n\n"
}

func Tab() string {
	return "\t"
}
func Line() string {
	return "\n---------------\n"
}
func DoubleLine() string {
	return "\n===============\n"
}
