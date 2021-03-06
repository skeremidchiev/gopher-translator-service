package configurate

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Config interface {
	GetVowelLetters() []string
	GetConsonantLetters() []string
	GetConsonantSounds() map[string][]string
	IsVowelLetter(string) bool
	IsConsonantLetter(string) bool
	CheckConsonantSounds(string, string) (string, error)
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

func (jc JsonConfig) CheckConsonantSounds(word, phonetics string) (string, error) {
	unknownError := fmt.Errorf("Unknown Consonant Sounds")

	// fmt.Printf("word: %+v\n", word)
	// fmt.Printf("phonetics: %+v\n", phonetics)

	for key, val := range jc.ConsonantSounds {
		if strings.HasPrefix(phonetics, key) {

			// fmt.Printf("sound: %+v\n", key)
			// fmt.Printf("digraphs: %+v\n", val)

			for _, v := range val {
				// fmt.Printf("Prefix: %+v\n", v)
				if strings.HasPrefix(word, v) {
					// fmt.Printf("Prefix: %+v  DONE\n", v)
					return v, nil
				}
			}
		}
	}

	return "", unknownError
}

func (jc JsonConfig) GetDigraphsFromSound(phonetics string) []string {
	return jc.ConsonantSounds[phonetics]
}

func (jc JsonConfig) IsSoundMatchingDigraph(phonetics string, digraph string) bool {
	digraphs := jc.ConsonantSounds[phonetics]
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
