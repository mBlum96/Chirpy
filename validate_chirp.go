package main

import (
	"encoding/json"
	"log"
	"net/http"
)

const MAX_CHIRP = 140

type ErrType int

const (
	ErrDecoding ErrType = iota
	ErrTooLong
)

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
		handleErr(w, r, ErrDecoding)
		return
	}
	if len(params.Body) > MAX_CHIRP {
		handleErr(w, r, ErrTooLong)
		return
	}
	handleSuccess(w, r)
}

func handleErr(w http.ResponseWriter, r *http.Request, errNumber ErrType) {
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
	return
}

func handleSuccess(w http.ResponseWriter, r *http.Request) {
	type returnVals struct {
		Valid bool `json:"valid"`
	}
	respBody := returnVals{
		Valid: true,
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
	return
}
