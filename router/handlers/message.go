package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/unity-web-service/messages"
)

func UpsertMessage(m messages.Repo) http.HandlerFunc {
	var message messages.Message

	return func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		if err := json.Unmarshal(b, &message); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		// validate the data in the struct is correct
		if len(message.Msg) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("error: please add at least one message to your request"))
			return
		}
		// checking correct format, still need to check if the ip address is valid
		if net.ParseIP(message.SentFromIP) == nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("error: please use proper IPv4 format: x.x.x.x"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("%v", message)))
	}
}
