package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

const MAX_CHIRP = 140

type ErrType int

const (
	ErrDecoding ErrType = iota
	ErrTooLong
)

type Profanity string

const (
	kerfuffle Profanity = "kerfuffle"
	sharbert  Profanity = "sharbert"
	fornax    Profanity = "fornax"
)

const profanityReplacement string = "****"

var profanities = []Profanity{kerfuffle, sharbert, fornax}

const ErrDecodingHeader int = 500
const ErrTooLongHeader int = 400

func (cfg *apiConfig) handleValidation(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondErr(w, ErrDecoding)
		return
	}
	if len(params.Body) > MAX_CHIRP {
		respondErr(w, ErrTooLong)
		return
	}
	// params.Body = cleanBody(params.Body)
	respondJSON(w, cleanBody(params.Body))
}

func respondErr(w http.ResponseWriter, errNumber ErrType) {
	var errHeader int
	type returnVals struct {
		Error string `json:"error"`
	}
	var respBody returnVals
	if errNumber == ErrDecoding {
		errHeader = ErrDecodingHeader
		respBody = returnVals{
			Error: "Something went wrong",
		}
	} else if errNumber == ErrTooLong {
		errHeader = ErrTooLongHeader
		respBody = returnVals{
			Error: "Chirp is too long",
		}
	}
	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errHeader)
	w.Write(dat)
}

func respondJSON(w http.ResponseWriter, cleanedBody string) {
	type returnVals struct {
		Body string `json:"cleaned_body"`
	}
	respBody := returnVals{
		Body: cleanedBody,
	}
	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)
}

func cleanBody(body string) (cleanBody string) {
	words := strings.Fields(body)
	for i, word := range words {
		for _, profanity := range profanities {
			if strings.Contains(strings.ToLower(word), string(profanity)) {
				words[i] = profanityReplacement
				break
			}
		}
	}
	return strings.Join(words, " ")
}
