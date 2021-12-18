package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // import postgres
	"github.com/shaural/GlobePersonalWebsite/server/pkg/common"
)

// Database ...
type Database interface {
	Initialize() error
	Close() error
	InsertCountry(*Country) error
	InsertState(*State) error
	UpdateCard(*Card) error
	GetCountries() ([]*Country, error)
	GetStates(string) ([]*State, error)
	GetCards() ([]*Card, error)
}

type gormDb struct {
	Database *gorm.DB
	ctx      context.Context
}

// NewDatabase initializes a new instance of a Database
func NewDatabase(ctx context.Context) (Database, error) {
	config := common.Config()
	db, err := gorm.Open("postgres", fmt.Sprintf("%s?sslmode=disable", config.DatabaseURL))
	if err != nil {
		return nil, err
	}
	return &gormDb{
		Database: db,
		ctx:      ctx,
	}, nil
}

func (gdb *gormDb) Initialize() error {
	if err := gdb.Database.Transaction(func(tx *gorm.DB) error {
		log.Println("gorm: Automigrate tables")
		return tx.AutoMigrate(&Country{}, &State{}, &Card{}).Error
	}); err != nil {
		return err
	}
	return nil
}

func (gdb *gormDb) Close() error {
	return gdb.Database.Close()
}

func (gdb *gormDb) InsertCountry(country *Country) error {
	return gdb.Database.Create(&country).Error
}

func (gdb *gormDb) InsertState(state *State) error {
	return gdb.Database.Create(&state).Error
}

func (gdb *gormDb) UpdateCard(card *Card) error {
	return gdb.Database.
		Where(Card{CountryID: card.CountryID, Title: card.Title}).
		Assign(&Card{
			StateID:       card.StateID,
			Title:         card.Title,
			Description:   card.Description,
			StartDate:     card.StartDate,
			EndDate:       card.EndDate,
			ImgFolderPath: card.ImgFolderPath,
			Link:          card.Link,
			Github:        card.Github,
			Type:          card.Type,
		}).
		FirstOrCreate(&card).Error
}

func (gdb *gormDb) GetCountries() (countries []*Country, err error) {
	return countries, gdb.Database.Find(&countries).Error
}

func (gdb *gormDb) GetStates(country string) (states []*State, err error) {
	gormDB := gdb.Database
	if len(country) > 0 {
		gormDB = gormDB.Where("country_id = ?", country)
	}
	return states, gormDB.Find(&states).Error
}

func (gdb *gormDb) GetCards() (cards []*Card, err error) {
	return cards, gdb.Database.Find(&cards).Error
}
