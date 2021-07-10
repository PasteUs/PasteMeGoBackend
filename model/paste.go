package model

type Paste interface {
	Save() error
	Get() error
	Delete() error
	GetContent() string
	GetLang() string
	GetPassword() string
}
