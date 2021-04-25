package translater

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/skeremidchiev/gopher-translator-service/app/configurate"
	"github.com/skeremidchiev/gopher-translator-service/app/lib/rateLimitter"
)

const (
	dictionaryapi_us = "https://api.dictionaryapi.dev/api/v2/entries/en/%s"
)

type Translater struct {
	conf configurate.Config
	rl   *rateLimitter.Limitter
}

func NewTranslater(c configurate.Config) *Translater {

	return &Translater{
		conf: c,
		rl:   rateLimitter.NewRateLimitter(1, 100*time.Millisecond),
	}
}

func (t *Translater) ParseWord(word string) (string, error) {
	word = strings.ToLower(word)

	if checkForEmptyWord(word) {
		return "", fmt.Errorf("[Translater] Word is not provided!")
	}

	if checkForApostrophes(word) {
		return "", fmt.Errorf("[Translater] Word \"%s\" contains apostrophes!", word)
	}

	if t.conf.IsVowelLetter(string(word[0])) {
		return "g" + word, nil
	}

	if word[:2] == "xr" {
		return "ge" + word, nil
	}

	ph, err := t.getPhonetics(word)
	if err != nil {
		return "", err
	}

	prefix, err := t.conf.CheckConsonantSounds(word, ph)
	if err != nil {
		return "", err
	}

	wordWithoutPrefix := word[len(prefix):]
	if strings.HasPrefix(wordWithoutPrefix, "qu") {
		return wordWithoutPrefix[2:] + prefix + "qu" + "ogo", nil
	}

	return wordWithoutPrefix + prefix + "ogo", nil
}

func (t *Translater) ParseSentence(words string) (string, error) {

	return "", nil
}

// getPhonetics - retrieves word phonetics from external API in order
// to find whether word starts with consonant sound
func (t *Translater) getPhonetics(word string) (string, error) {
	t.rl.Throttle()

	url := fmt.Sprintf(dictionaryapi_us, word)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	j := []struct {
		Word      string `json="word"`
		Phonetics []struct {
			Text string `json="text"`
		} `json="phonetics"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&j)
	if err != nil {
		return "", err
	}

	phonetics := j[0].Phonetics[0].Text
	if len(phonetics) < 3 { // format is "/sounds/"
		return "", fmt.Errorf("[Translater] Incorrect or unknown word")
	}

	phonetics = strings.ReplaceAll(phonetics, "/", "") // remove forward slashes

	return phonetics, nil
}

func checkForEmptyWord(word string) bool {
	if len(word) == 0 {
		return true
	}

	return false
}

func checkForApostrophes(word string) bool {
	if strings.ContainsAny(word, "’'") {
		return true
	}

	return false
}
