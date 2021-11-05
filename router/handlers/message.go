package handlers

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"

	"github.com/unity-web-service/messages"
	"github.com/unity-web-service/queue"
)

var (
	ERR_INTERNAL = errors.New("internal server error: couldn't unmarshal request into map")
)

func UpsertMessage(m messages.Repo, p queue.IPublisher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var message messages.MessageRequest

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		if err := checkRequestKeys(b); err != nil {
			if err == ERR_INTERNAL {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		if err := json.Unmarshal(b, &message); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		if err := validatemessagereq(message); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		messageToSendToDBAndQ, err := convertType(message)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		if err := m.UpsertMessage(r.Context(), messageToSendToDBAndQ); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("message added"))
	}
}

// if priority isn't set, it will default to 0
func validatemessagereq(m messages.MessageRequest) error {
	if m.Sender == nil || len(*m.Sender) == 0 {
		return errors.New("please add a sender")
	}
	if m.Ts == nil {
		return errors.New("please add a valid unix timestamp")
	}
	if _, err := strconv.ParseInt(*m.Ts, 10, 64); err != nil {
		return errors.New("not a valid unix timestamp")
	}
	if len(m.Msg) == 0 {
		return errors.New("please add at least one message to your request")
	}
	if m.SentFromIP == nil {
		return errors.New("please add a valid ip address")
	}
	if net.ParseIP(*m.SentFromIP) == nil {
		return errors.New("please use proper IPv4 format: xxx.xxx.xx.xx")
	}

	return nil
}

func convertType(m messages.MessageRequest) (*messages.Message, error) {
	var messageToSendToDBAndQ messages.Message

	b, err := json.Marshal(m.Msg)
	if err != nil {
		return nil, errors.New("unable to convert messages into byte slice")
	}
	i, err := strconv.Atoi(*m.Ts)
	if err != nil {
		return nil, errors.New("unable to convert timestamp to int")
	}

	messageToSendToDBAndQ.Msg = hex.EncodeToString(b)
	messageToSendToDBAndQ.Timestamp = i
	messageToSendToDBAndQ.Priority = m.Priority
	messageToSendToDBAndQ.Sender = *m.Sender
	messageToSendToDBAndQ.SentFromIP = *m.SentFromIP

	return &messageToSendToDBAndQ, nil
}

// this helper is going to check the keys on the request
// and error if there is a key that isn't supposed to be
//on the request and return nil if the keys are correct
func checkRequestKeys(b []byte) error {
	requestKeyMap := make(map[string]interface{})
	validKeys := []string{"ts", "sender", "sent-from-ip", "message", "priority"}

	if err := json.Unmarshal(b, &requestKeyMap); err != nil {
		return ERR_INTERNAL
	}

	for _, k := range validKeys {
		delete(requestKeyMap, k)
	}

	if len(requestKeyMap) > 0 {
		keys := make([]string, 0)
		for k, _ := range requestKeyMap {
			keys = append(keys, k)
		}
		return fmt.Errorf("invalid keys in request, please remove them: %v", keys)
	}

	return nil
}
