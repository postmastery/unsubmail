package main

import (
	"fmt"
	"regexp"
)

// replace is similar to regexp.ReplaceAllString, but instead of returning a copy of src it returns the template
func replace(re *regexp.Regexp, template, src string) (string, error) {

	match := re.FindStringSubmatchIndex(src)
	if match == nil {
		return "", fmt.Errorf("no match found: %q", src)
	}
	//log.Printf("match: %v", match)

	var dst []byte
	dst = re.ExpandString(dst, template, src, match)

	return string(dst), nil
}
