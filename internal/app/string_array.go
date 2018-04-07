package app

import (
	"strings"
)

type StringArray []string

// Support comma separator flag for multiple times case
func (a *StringArray) Set(s string) error {
	if strings.Contains(s, ",") {
		b := strings.Split(s, ",")
		for i := 0; i < len(b); i++ {
			*a = append(*a, b[i])
		}
	} else {
		*a = append(*a, s)
	}
	return nil
}

func (a *StringArray) String() string {
	return strings.Join(*a, ",")
}
