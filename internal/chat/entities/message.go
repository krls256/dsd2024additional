package entities

import "encoding/json"

type Message struct {
	Nickname string `json:"nickname"`
	Message  string `json:"message"`
}

func (m Message) MarshalBinary() (data []byte, err error) {
	return json.Marshal(m)
}
