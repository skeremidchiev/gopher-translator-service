package storage

import (
	"fmt"
	"sort"
	"strings"
)

type Storage interface {
	Exists(string) bool
	Save(string, string) error
	Update(string, string) error
	GetTranslation(string) (string, error)
	GetAll() string
}

type MapStorage struct {
	History map[string]string `json:"history" bson:"history"`
}

func (ms *MapStorage) Exists(key string) bool {
	if _, ok := ms.History[key]; ok {
		return true
	}
	return false
}

func (ms *MapStorage) Save(key string, value string) error {
	if ms.Exists(key) {
		return fmt.Errorf("[Storage] Key already exists!")
	}

	ms.History[key] = value
	return nil
}

func (ms *MapStorage) Update(key string, value string) error {
	ms.History[key] = value
	return nil
}

func (ms *MapStorage) GetTranslation(key string) (string, error) {
	if !ms.Exists(key) {
		return "", fmt.Errorf("[Storage] Key doesn't exists!")
	}

	return ms.History[key], nil
}

type sortableSlice []string

func (ss sortableSlice) Len() int {
	return len(ss)
}
func (ss sortableSlice) Swap(i, j int) {
	ss[i], ss[j] = ss[j], ss[i]
}
func (ss sortableSlice) Less(i, j int) bool {
	return strings.Compare(ss[i], ss[j]) == -1
}

func (ms *MapStorage) GetAll() string {
	keys := make(sortableSlice, 0)

	for key := range ms.History {
		keys = append(keys, key)
	}

	sort.Sort(keys)

	result := "["
	for i, key := range keys {
		if i != 0 {
			result += ","
		}
		result += fmt.Sprintf("{\"%s\":\"%s\"}", key, ms.History[key])
	}
	result += "]"

	return result
}

func NewStorage() Storage {
	return &MapStorage{
		History: make(map[string]string),
	}
}
