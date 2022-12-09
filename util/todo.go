package util

import (
	"github.com/prometheus/common/log"
	"net/http"
)

func TodoEvent(w http.ResponseWriter) {
	_, err := w.Write([]byte{})
	if err != nil {
		log.Errorln(err)
	}
}
