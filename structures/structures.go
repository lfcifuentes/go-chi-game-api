package structures

import (
	"errors"
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type User struct {
	Id        uuid.UUID
	Username  string
	Scores    []Score
	CreatedAt time.Time
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.Id, err = uuid.NewV1()

	if err != nil {
		return errors.New("can't save invalid data")
	}
	return
}

func (u *User) IsValid() bool {
	if u.Id.String() == "00000000-0000-0000-0000-000000000000" {
		return false
	}
	return true
}

type Score struct {
	Id        uuid.UUID
	User_id   uuid.UUID
	Score     int
	CreatedAt time.Time
}

type BestScores struct {
	Score     int
	Username  string
	CreatedAt time.Time
}

func (s *Score) BeforeCreate(tx *gorm.DB) (err error) {
	s.Id, err = uuid.NewV1()

	if err != nil {
		return errors.New("can't save invalid data")
	}
	return
}

type ResponseUser struct {
	Status  int    `json:"status"`
	Data    User   `json:"data"`
	Message string `json:"message"`
}

type ResponseScores struct {
	Status  int          `json:"status"`
	Data    []BestScores `json:"data"`
	Message string       `json:"message"`
}
