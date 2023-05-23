package language

import (
	"sync"
)

type lang struct {
	lock sync.RWMutex

	data map[string]string
}

func (a *lang) SetLanguage(source string, target string) *lang {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.data[source] = target
	return a
}

func (a *lang) GetLanguage(source string) string {
	a.lock.Lock()
	defer a.lock.Unlock()
	s, b := a.data[source]
	if !b {
		return source
	}
	return s
}

var Language = lang{data: map[string]string{}}

func L(source string) string {
	return Language.GetLanguage(source)
}
