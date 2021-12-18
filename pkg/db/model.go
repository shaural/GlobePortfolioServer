package db

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Country model for country metadata
type Country struct {
	ID string
	Name string
	Latitude int
	Longitude int
}

// State model for state metadata
type State struct {
	ID string
	CountryID string
	Name string
}

// Card model for storing any information (images, projects, text, education)
type Card struct {
	gorm.Model
	CountryID string
	StateID string
	Title string
	Description string
	StartDate time.Time
	EndDate time.Time
	ImgFolderPath string
	Link string
	Github string
	Type string
}
