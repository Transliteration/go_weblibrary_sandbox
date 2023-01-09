package main

import (
	"github.com/google/uuid"
)

// Пользователь
// - имя
// - uuid v4
// - возраст (опционально)
// - список книг

type User struct {
	name   string
	userID uuid.UUID
	age    int
	books  []uuid.UUID
}

// Книга
// - название
// - жанр (опционально)

type Genre int64

const (
	Thriller Genre = iota
	Comedy
	Drama
)

type Book struct {
	bookID uuid.UUID
	name   string
	genre  []Genre
}
