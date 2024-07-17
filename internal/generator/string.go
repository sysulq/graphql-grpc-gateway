package generator

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func Title(name string) string {
	title := cases.Title(language.English, cases.NoLower)
	return title.String(name)
}
