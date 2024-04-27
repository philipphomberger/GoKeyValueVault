package models

type Secret struct {
	Key   string
	Value string `json:",omitempty"`
}
