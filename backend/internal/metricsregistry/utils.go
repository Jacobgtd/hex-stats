package metricsregistry

import "strings"

func templateParams(template string) []string {
	var params []string

	for _, part := range strings.Split(strings.Trim(template, "/"), "/") {
		if strings.HasPrefix(part, ":") {
			params = append(params, part[1:])
		}
	}

	return params
}
