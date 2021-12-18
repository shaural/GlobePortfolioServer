package db

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/suite"
)

// AnyTime type created to mock time
type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	return true
}

type gormFixture struct {
	suite.Suite
	gormdb *gormDb
	mockdb *sql.DB
	mock   sqlmock.Sqlmock
}

func (s *gormFixture) BeforeTest(suiteName, testName string) {
	mockdb, mock, err := sqlmock.New()
	s.Nil(err)

	s.mockdb = mockdb
	s.mock = mock
	gdb, err := gorm.Open("postgres", mockdb)
	s.Nil(err)

	s.gormdb = &gormDb{
		Database: gdb,
		ctx:      context.Background(),
	}
}

func (s *gormFixture) AfterTest(suiteName, testName string) {
	s.gormdb.Database.Close()
	s.Nil(s.mock.ExpectationsWereMet())
}

type countryFixture struct {
	gormFixture

	rows    *sqlmock.Rows
	country *Country
}

func TestCountryFixture(t *testing.T) {
	suite.Run(t, new(countryFixture))
}

func (f *countryFixture) SetupTest() {
	f.rows = sqlmock.NewRows([]string{"id", "name", "latitude", "longitude"}).
		AddRow("APPLE", "BANANA", 12, 34)
	f.country = &Country{
		ID:        "APPLE",
		Name:      "BANANA",
		Latitude:  12,
		Longitude: 34,
	}
}

func (f *countryFixture) Test_GetCountries() {
	f.mock.ExpectQuery("^SELECT.*").WillReturnRows(f.rows)

	res, err := f.gormdb.GetCountries()

	f.Equal(1, len(res))
	f.Nil(err)
	f.Equal(f.country, res[0])
}

func (f *countryFixture) Test_GetCountries_error() {
	f.mock.ExpectQuery("^SELECT.*").WillReturnError(fmt.Errorf("QUERY ERROR"))

	res, err := f.gormdb.GetCountries()

	f.EqualError(err, "QUERY ERROR")
	f.Equal([]*Country{}, res)
}

func (f *countryFixture) Test_InsertCountry() {
	retIDRow := sqlmock.NewRows([]string{"id"}).
		AddRow(f.country.ID)

	f.mock.ExpectBegin()
	f.mock.ExpectQuery("^INSERT.*countries.*").
		WithArgs(f.country.ID, f.country.Name, f.country.Latitude, f.country.Longitude).
		WillReturnRows(retIDRow)
	f.mock.ExpectCommit()

	err := f.gormdb.InsertCountry(f.country)

	f.Nil(err)
}

func (f *countryFixture) Test_InsertCountry_Error() {
	f.mock.ExpectBegin()
	f.mock.ExpectQuery("^INSERT.*countries.*").
		WithArgs(f.country.ID, f.country.Name, f.country.Latitude, f.country.Longitude).
		WillReturnError(fmt.Errorf("INSERT ERROR"))
	f.mock.ExpectRollback()

	err := f.gormdb.InsertCountry(f.country)

	f.EqualError(err, "INSERT ERROR")
}

type stateFixture struct {
	gormFixture

	rows  *sqlmock.Rows
	state *State
}

func TestStateFixture(t *testing.T) {
	suite.Run(t, new(stateFixture))
}

func (f *stateFixture) SetupTest() {
	f.state = &State{
		ID:        "APPLE",
		CountryID: "BANANA",
		Name:      "KIWI",
	}
	f.rows = sqlmock.NewRows([]string{"id", "country_id", "name"}).
		AddRow(f.state.ID, f.state.CountryID, f.state.Name)
}

func (f *stateFixture) Test_GetStates() {
	f.mock.ExpectQuery("^SELECT.*").WillReturnRows(f.rows)

	res, err := f.gormdb.GetStates("")

	f.Equal(1, len(res))
	f.Nil(err)
	f.Equal(f.state, res[0])
}

func (f *stateFixture) Test_GetStates_error() {
	f.mock.ExpectQuery("^SELECT.*").WillReturnError(fmt.Errorf("QUERY ERROR"))

	res, err := f.gormdb.GetStates("")

	f.EqualError(err, "QUERY ERROR")
	f.Equal([]*State{}, res)
}

func (f *stateFixture) Test_GetStates_country() {
	f.mock.ExpectQuery("^SELECT.*WHERE").WillReturnRows(f.rows)

	res, err := f.gormdb.GetStates("US")

	f.Equal(1, len(res))
	f.Nil(err)
	f.Equal(f.state, res[0])
}

func (f *stateFixture) Test_InsertState() {
	retIDRow := sqlmock.NewRows([]string{"id"}).
		AddRow(f.state.ID)

	f.mock.ExpectBegin()
	f.mock.ExpectQuery("^INSERT.*states.*").
		WithArgs(f.state.ID, f.state.CountryID, f.state.Name).
		WillReturnRows(retIDRow)
	f.mock.ExpectCommit()

	err := f.gormdb.InsertState(f.state)

	f.Nil(err)
}

func (f *stateFixture) Test_InsertState_Error() {
	f.mock.ExpectBegin()
	f.mock.ExpectQuery("^INSERT.*states.*").
		WithArgs(f.state.ID, f.state.CountryID, f.state.Name).
		WillReturnError(fmt.Errorf("INSERT ERROR"))
	f.mock.ExpectRollback()

	err := f.gormdb.InsertState(f.state)

	f.EqualError(err, "INSERT ERROR")
}

type cardFixture struct {
	gormFixture

	rows *sqlmock.Rows
	t    time.Time
	card *Card
}

func TestCardFixture(t *testing.T) {
	suite.Run(t, new(cardFixture))
}

func (f *cardFixture) SetupTest() {
	f.t = time.Now()
	f.rows = sqlmock.NewRows([]string{
		"country_id",
		"state_id",
		"title",
		"description",
		"start_date",
		"end_date",
		"img_folder_path",
		"link",
		"github",
		"type"}).
		AddRow("APPLE",
			"BANANA",
			"KIWI",
			"MANGO",
			f.t,
			f.t,
			"ORANGE",
			"GRAPE",
			"BERRY",
			"PEAR")
	f.card = &Card{
		CountryID:     "APPLE",
		StateID:       "BANANA",
		Title:         "KIWI",
		Description:   "MANGO",
		StartDate:     f.t,
		EndDate:       f.t,
		ImgFolderPath: "ORANGE",
		Link:          "GRAPE",
		Github:        "BERRY",
		Type:          "PEAR",
	}
}

func (f *cardFixture) Test_GetCards() {
	f.mock.ExpectQuery("^SELECT.*").WillReturnRows(f.rows)

	res, err := f.gormdb.GetCards()

	f.Equal(1, len(res))
	f.Nil(err)
	f.Equal(f.card, res[0])
}

func (f *cardFixture) Test_GetCards_error() {
	f.mock.ExpectQuery("^SELECT.*").WillReturnError(fmt.Errorf("QUERY ERROR"))

	res, err := f.gormdb.GetCards()

	f.EqualError(err, "QUERY ERROR")
	f.Equal([]*Card{}, res)
}

func (f *cardFixture) Test_UpdateCard_Update() {
	res := sqlmock.NewResult(123, 456)
	f.mock.ExpectQuery("^SELECT.*cards.*").
		WithArgs(f.card.CountryID, f.card.Title).
		WillReturnRows(f.rows)
	f.mock.ExpectBegin()
	f.mock.ExpectExec("^UPDATE.*cards.*").
		WithArgs("UPDATED",
			f.card.EndDate,
			f.card.Github,
			f.card.ImgFolderPath,
			f.card.Link,
			f.card.StartDate,
			f.card.StateID,
			f.card.Title,
			f.card.Type,
			AnyTime{},
			f.card.CountryID,
			f.card.Title).
		WillReturnResult(res)
	f.mock.ExpectCommit()

	err := f.gormdb.UpdateCard(&Card{
		CountryID:     "APPLE",
		StateID:       "BANANA",
		Title:         "KIWI",
		Description:   "UPDATED",
		StartDate:     f.t,
		EndDate:       f.t,
		ImgFolderPath: "ORANGE",
		Link:          "GRAPE",
		Github:        "BERRY",
		Type:          "PEAR",
	})

	f.Nil(err)
}

func (f *cardFixture) Test_UpdateCard_Insert() {
	emptyRow := sqlmock.NewRows([]string{
		"country_id",
		"state_id",
		"title",
		"description",
		"start_date",
		"end_date",
		"img_folder_path",
		"link",
		"github",
		"type"})
	retIDRow := sqlmock.NewRows([]string{"id"}).
		AddRow(f.card.ID)
	f.mock.ExpectQuery("^SELECT.*cards.*").
		WithArgs(f.card.CountryID, f.card.Title).
		WillReturnRows(emptyRow)
	f.mock.ExpectBegin()
	f.mock.ExpectQuery("^INSERT INTO .*cards.*").
		WithArgs(AnyTime{},
			AnyTime{},
			AnyTime{},
			f.card.CountryID,
			f.card.StateID,
			f.card.Title,
			f.card.Description,
			f.card.StartDate,
			f.card.EndDate,
			f.card.ImgFolderPath,
			f.card.Link,
			f.card.Github,
			f.card.Type).
		WillReturnRows(retIDRow)
	f.mock.ExpectCommit()

	err := f.gormdb.UpdateCard(f.card)

	f.Nil(err)
}

func (f *cardFixture) Test_UpdateCard_SelectError() {
	f.mock.ExpectQuery("^SELECT.*cards.*").
		WithArgs(f.card.CountryID, f.card.Title).
		WillReturnError(fmt.Errorf("SELECT ERROR"))

	err := f.gormdb.UpdateCard(f.card)

	f.EqualError(err, "SELECT ERROR")
}

func (f *cardFixture) Test_UpdateCard_Update_Error() {
	f.mock.ExpectQuery("^SELECT.*cards.*").
		WithArgs(f.card.CountryID, f.card.Title).
		WillReturnRows(f.rows)
	f.mock.ExpectBegin()
	f.mock.ExpectExec("^UPDATE.*cards.*").
		WithArgs("UPDATED",
			f.card.EndDate,
			f.card.Github,
			f.card.ImgFolderPath,
			f.card.Link,
			f.card.StartDate,
			f.card.StateID,
			f.card.Title,
			f.card.Type,
			AnyTime{},
			f.card.CountryID,
			f.card.Title).
		WillReturnError(fmt.Errorf("UPDATE ERROR"))
	f.mock.ExpectRollback()

	err := f.gormdb.UpdateCard(&Card{
		CountryID:     "APPLE",
		StateID:       "BANANA",
		Title:         "KIWI",
		Description:   "UPDATED",
		StartDate:     f.t,
		EndDate:       f.t,
		ImgFolderPath: "ORANGE",
		Link:          "GRAPE",
		Github:        "BERRY",
		Type:          "PEAR",
	})

	f.EqualError(err, "UPDATE ERROR")
}

func (f *cardFixture) Test_UpdateCard_Insert_Error() {
	emptyRow := sqlmock.NewRows([]string{
		"country_id",
		"state_id",
		"title",
		"description",
		"start_date",
		"end_date",
		"img_folder_path",
		"link",
		"github",
		"type"})
	f.mock.ExpectQuery("^SELECT.*cards.*").
		WithArgs(f.card.CountryID, f.card.Title).
		WillReturnRows(emptyRow)
	f.mock.ExpectBegin()
	f.mock.ExpectQuery("^INSERT INTO .*cards.*").
		WithArgs(AnyTime{},
			AnyTime{},
			AnyTime{},
			f.card.CountryID,
			f.card.StateID,
			f.card.Title,
			f.card.Description,
			f.card.StartDate,
			f.card.EndDate,
			f.card.ImgFolderPath,
			f.card.Link,
			f.card.Github,
			f.card.Type).
		WillReturnError(fmt.Errorf("INSERT ERROR"))
	f.mock.ExpectRollback()

	err := f.gormdb.UpdateCard(f.card)

	f.EqualError(err, "INSERT ERROR")
}
