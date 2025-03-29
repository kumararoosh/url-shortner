package main

import (
	"github.com/google/uuid"
	_ "github.com/google/uuid"
)

func shorten_url(original_url string) string {
	id := uuid.New()
	return id.String()
}