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
)

func InsertMessage(m messages.Repo) http.HandlerFunc {
	var message messages.MessageRequest

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

		if err := m.InsertMessage(r.Context(), messageToSendToDBAndQ); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		//todo: send to db and nats here

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("%v", messageToSendToDBAndQ)))
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
