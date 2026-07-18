package pathutil

import "path"

func Split(p string) []string {
	dir, last := path.Split(p)
	if last == "" {
		return []string{}
	}
	return append(Split(path.Clean(dir)), last)
}
