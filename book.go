package main

import "database/sql"

//Model

// Book representation
type Book struct {
	ID     int
	Author string
	Name   string
	Pages  int
	ISBN   string
}

// BookRepository resource
type BookRepository struct {
	storage sql.DB
}

// Create a book
func (r *BookRepository) Create(b *Book) (*Book, error) {
	if err := r.storage.QueryRow(
		"INSERT INTO books (author, name, pages, isbn) VALUES ($1, $2, $3, $4) RETURNING id",
		b.Author, b.Name,
		b.Pages, b.ISBN,
	).Scan(&b.ID); err != nil {
		return nil, err
	}
	return b, nil

}

// FindByAuthor ...
func (r *BookRepository) FindByAuthor(author string) (*Book, error) {
	return nil, nil
}
