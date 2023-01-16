package models

import (
	"github.com/google/uuid"
)

// Пользователь
// - имя
// - uuid v4
// - возраст (опционально)
// - список книг

type User struct {
	UserID uuid.UUID
	Name   string
	Age    int
}

// Книга
// - название
// - жанр (опционально)

// type Genre int64

// const (
// 	Thriller Genre = iota
// 	Comedy
// 	Drama
// )

type Book struct {
	BookID uuid.UUID
	Name   string
	// genre  []Genre
}
