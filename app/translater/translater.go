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

	if t.conf.IsVowelLetter(string(word[0])) { // vowel case
		return "g" + word, nil
	}

	if word[:2] == "xr" { // "xr" case
		return "ge" + word, nil
	}

	ph, err := t.getPhonetics(word)
	if err != nil {
		return "", fmt.Errorf("[Translater] Word \"%s\" failed with error: %s!", word, err.Error())
	}

	fmt.Printf("ph: %+v\n", ph)

	prefix, err := t.conf.CheckConsonantSounds(word, ph)
	if err != nil {
		return "", fmt.Errorf("[Translater] Word \"%s\" failed with error: %s!", word, err.Error())
	}

	wordWithoutPrefix := word[len(prefix):]
	if strings.HasPrefix(wordWithoutPrefix, "qu") { // consonant sound and "qu"
		return wordWithoutPrefix[2:] + prefix + "qu" + "ogo", nil
	}

	return wordWithoutPrefix + prefix + "ogo", nil // consonant sound
}

func (t *Translater) ParseSentence(sentence string) (string, error) {
	if !checkForSign(sentence) {
		return "", fmt.Errorf("[Translater] Sentence \"%s\" badly formed!", sentence)
	}

	result := ""
	l := len(sentence)
	sign := sentence[l-1:]
	sentence = sentence[:l-1]
	for i, word := range strings.Split(sentence, " ") {
		translation, err := t.ParseWord(word)
		if err != nil {
			return "", err
		}

		if i == 0 { // first word to upper case
			translation = strings.Title(translation)
		}

		result += translation + " "
	}
	result = result[:len(result)-1] + sign

	return result, nil
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
		return "", fmt.Errorf("[Translater] Unknown word")
	}

	phonetics := j[0].Phonetics[0].Text
	if len(phonetics) < 3 { // format is "/sounds/"
		return "", fmt.Errorf("[Translater] Incorrect or unknown word")
	}

	phonetics = strings.ReplaceAll(phonetics, "/", "")   // remove forward slashes
	phonetics = strings.ReplaceAll(phonetics, "ˈ", "")   // remove emphasis
	phonetics = strings.ReplaceAll(phonetics, "ˌ", "")   // remove separators
	phonetics = strings.ReplaceAll(phonetics, "(h)", "") // remove silent h
	phonetics = strings.ReplaceAll(phonetics, "(n)", "") // remove silent n
	phonetics = strings.ReplaceAll(phonetics, "(p)", "") // remove silent p

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

func checkForSign(sentence string) bool {
	if strings.LastIndexAny(sentence, ".!?") != len(sentence)-1 {
		return false
	}

	return true
}
