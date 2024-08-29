package domain

import "encoding/json"

type NewFileEvent struct {
	FileName string `json:"file_name"`
}

func (n *NewFileEvent) ToJson() (string, error) {
	bytes, error := json.Marshal(n)

	return string(bytes), error
}

func NewNewFileEventFromJson(jsonString string) (*NewFileEvent, error) {
	n := &NewFileEvent{}

	err := json.Unmarshal([]byte(jsonString), n)

	return n, err
}
