package util

import (
	"log"
	"net/http"
)

func TodoEvent(w http.ResponseWriter) {
	_, err := w.Write([]byte{})
	if err != nil {
		log.Fatalln(err)
	}
}
