package model

type App struct {
	ID     int64
	Name   string
	Secret string // TODO: убрать, хранить в зашифрованном виде
}
