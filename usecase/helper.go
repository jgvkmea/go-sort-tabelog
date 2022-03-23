package usecase

import "strings"

func getSearchCondition(message string) []string {
	return strings.Split(strings.ReplaceAll(message, "　", " "), " ")
}
