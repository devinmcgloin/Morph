package search

import (
	"strings"
)

func formatQueryString(req []string, opt []string, exc []string) string {
	args := []string{}

	trim(req)
	req = filterEmpty(req)
	trim(opt)
	opt = filterEmpty(opt)
	trim(exc)
	exc = filterEmpty(exc)

	for i, ex := range exc {
		exc[i] = "!" + ex
	}
	if len(req) != 0 {
		args = append(args, "("+strings.Join(req, " & ")+")")
	}

	if len(opt) != 0 {
		args = append(args, "("+strings.Join(opt, " | ")+")")
	}

	if len(exc) != 0 {
		args = append(args, strings.Join(exc, " & "))
	}
	return strings.Join(args, " & ")
}

func filter(vs []string, f func(string) bool) []string {
	vsf := make([]string, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func trim(vs []string) {
	for i, v := range vs {
		vs[i] = strings.Trim(v, ",./!@#$%^&*()_+-= ")
	}
}

func filterEmpty(arr []string) []string {
	empty := func(s string) bool {
		return s != ""
	}

	return filter(arr, empty)
}
