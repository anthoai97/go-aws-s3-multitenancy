package core

type ResponseChan[T any] struct {
	Data  T
	Error error
}
