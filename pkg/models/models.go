package models

import (
	"errors"
)

var ErrNoRecord = errors.New("models: no matching record found")

type Website struct {
	ID    int
	Title string
}

type Sub struct {
	ID        int
	Title     string
	WebsiteId int
}
