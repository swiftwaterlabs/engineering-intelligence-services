package core

import "strings"

func SanitizeString(value string) string {
	return strings.ToValidUTF8(value, " ")
}
