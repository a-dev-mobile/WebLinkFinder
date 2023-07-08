package regexutils

import (
	"fmt"
	"regexp"
)

func CompileRegexes(regexStrs []string) ([]*regexp.Regexp, error) {
	regexes := make([]*regexp.Regexp, len(regexStrs))

	for i, regexStr := range regexStrs {
		regex, err := regexp.Compile(regexStr)
		if err != nil {
			return nil, fmt.Errorf("error compiling regex \"%s\": %v", regexStr, err)
		}
		regexes[i] = regex
	}

	return regexes, nil
}
