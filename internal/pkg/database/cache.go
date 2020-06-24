package database

import (
	"strings"
	"sync"
)

var (
	ForumSlug map[string]string
	FMux sync.Mutex
)

func init()  {
	ForumSlug = make(map[string]string)

}
func ForumCheckSlug(slug string) (rslug string, state bool) {
	FMux.Lock()
	defer FMux.Unlock()
	i, err := ForumSlug[strings.ToLower(slug)]
	if err != false {
		return slug, false
	}
	return i, true
}

func ForumSetSlug(slug string) {
	FMux.Lock()
	defer FMux.Unlock()
	ForumSlug[strings.ToLower(slug)] = slug
}

func ForumClearSlug() bool {
	FMux.Lock()
	defer FMux.Unlock()
	ForumSlug = make(map[string]string)
	return false
}