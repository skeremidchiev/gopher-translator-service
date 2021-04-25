package configurate

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config interface {
	GetVowelLetters() []string
	GetConsonantLetters() []string
	GetConsonantSounds() map[string][]string
	IsVowelLetter(string) bool
	IsConsonantLetter(string) bool
	GetDigraphsFromSound(string) []string
	IsSoundMatchingDigraph(string, string) bool
}

type JsonConfig struct {
	VowelLetters     []string            `json:"vowelLetters" bson:"vowelLetters"`
	ConsonantLetters []string            `json:"consonantLetters" bson:"consonantLetters"`
	ConsonantSounds  map[string][]string `json:"consonantSounds" bson:"consonantSounds"`
}

func (jc JsonConfig) GetVowelLetters() []string {
	return jc.VowelLetters
}

func (jc JsonConfig) GetConsonantLetters() []string {
	return jc.ConsonantLetters
}

func (jc JsonConfig) GetConsonantSounds() map[string][]string {
	return jc.ConsonantSounds
}

func (jc JsonConfig) IsVowelLetter(letter string) bool {
	for _, v := range jc.VowelLetters {
		if v == letter {
			return true
		}
	}
	return false
}

func (jc JsonConfig) IsConsonantLetter(letter string) bool {
	for _, v := range jc.ConsonantLetters {
		if v == letter {
			return true
		}
	}
	return false
}

func (jc JsonConfig) GetDigraphsFromSound(sound string) []string {
	return jc.ConsonantSounds[sound]
}

func (jc JsonConfig) IsSoundMatchingDigraph(sound string, digraph string) bool {
	digraphs := jc.ConsonantSounds[sound]
	for _, d := range digraphs {
		if d == digraph {
			return true
		}
	}

	return false
}

func NewConfig(pathToConfigFile string) (Config, error) {
	jsonFile, err := os.Open(pathToConfigFile)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	config := &JsonConfig{}
	err = json.Unmarshal(byteValue, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
