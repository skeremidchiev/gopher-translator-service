# gopher-translator-service

## Project Overview
[**Free Dictionary API**](https://dictionaryapi.dev/) based service that parses sentences and words based on this rules:
- If a word starts with a vowel letter, add prefix “g” to the word (ex. apple => gapple)
- If a word starts with the consonant letters “xr”, add the prefix “ge” to the begging of the word. Such words as “xray” actually sound in the beginning with vowel sound as you pronounce them so a true gopher would say “gexray”.
- If a word starts with a consonant sound, move it to the end of the word and then add “ogo” suffix to the word. Consonant sounds can be made up of multiple consonants, a.k.a. a consonant cluster (e.g. "chair" -> "airchogo”).
- If a word starts with a consonant sound followed by "qu", move it to the end of the word, and then add "ogo" suffix to the word (e.g. "square" -> "aresquogo").

## Endpoints
- /word
- /sentence
- /history

## How to run
```
chmod u+x start.sh // give permissions once
./start.sh
```

## Limitations
Service doesn't handle commas in the sentences and fails when passed word is unknown for [**Free Dictionary API**](https://dictionaryapi.dev/)!