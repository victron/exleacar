package details

import (
	"errors"
	"regexp"
)

// finding link (https or http) in string
func FindLink(s string) (string, error) {
	var regLink, _ = regexp.Compile(`(https|http):\/\/[a-zA-Z0-9\/._\-%?=&]+`)
	result := regLink.FindString(s)
	if result == "" {
		return "", errors.New("no match")
	}
	return result, nil
}
