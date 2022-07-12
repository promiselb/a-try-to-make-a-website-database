package website

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Server struct {
	*mux.Router
	datas Slice[Data]
}

type Name string

type User struct {
	Firstname       Name `json:"firstname"`
	Lastname        Name `json:"lastname"`
	Nation          Name `json:"nation"`
	Birthday        `json:"birthday"`
	Gendre          `json:"gendre"`
	VIPSubscription `json:"VIPSubscription"`
	Gmail           `json:"gmail"`
	Password        `json:"password"`
	PhotoProfile    string    `json:"photoProfile"`
	ID              uuid.UUID `json:"uuid"`
}

type Day uint // 1 <= Day <= 31

type Month time.Month // 1 <= Month <= 12

type Year uint

type Gmail string

type Gendre string

type Password string

type Birthday struct {
	Day
	Month
	Year
}

type VIPSubscription struct {
	lasts        time.Duration
	subscription bool
}

type Slice[T any] []T

type Data struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// Check() checks if g contains at the end @gmail.com.
func (g Gmail) Check() (err error) {
	exists := strings.HasSuffix(string(g), "@gmail.com")
	if !exists {
		return errors.New("invalid gmail? Must contains '@gmail.com'! ")
	}
	return
}

// Check() checks if 1 <= d <= 31.
func (d Day) Check() (err error) {
	if 1 <= d {
		if d <= 31 {
			return
		} else {
			return fmt.Errorf("Day > 31. %d", d)
		}
	} else {
		return fmt.Errorf("Day < 1. %d", d)
	}
}

// Check() checks if 1 <= m <= 12.
func (m Month) Check() (err error) {
	if lessThan1, greaterThan12 := m < 1, m > 12; lessThan1 || greaterThan12 {
		return errors.New("invalid month of birth? Must be 1 <= m <= 12! ")
	}
	return
}

// Check() checks if 1 <= y.
func (y Year) Check() (err error) {
	if lessThan1 := y < 1; lessThan1 {
		return errors.New("invalid year of birth? Must be 1 < y! ")
	}
	return
}

// Check() checks if d, m and y are valid.
func (b Birthday) Check() (err error) {
	if ok := b.Day.Check() != nil; ok {
		return b.Day.Check()
	} else if ok := b.Month.Check() != nil; ok {
		return b.Month.Check()
	} else if ok := b.Year.Check() != nil; ok {
		return b.Year.Check()
	} else {
		return
	}
}

// Check() checks if the gendre is male or female.
func (g Gendre) Check() (err error) {
	s := strings.ToLower(string(g))
	if s != "male" && s != "female" {
		return errors.New("invalid Gendre?")
	}
	return
}

// Check() checks if the name is empty.
func (n Name) Check() (err error) {
	if n == "" {
		return errors.New("empty name?")
	}
	return
}

// Check() checks if the password is empty or contains illegal charaters.
func (p Password) Check() (err error) {
	if p == "" {
		return errors.New("empty password?")
	}

	if strings.Contains(string(p), "(") || strings.Contains(string(p), ")") || strings.Contains(string(p), "{") || strings.Contains(string(p), "}") || strings.Contains(string(p), "[") || strings.Contains(string(p), "]") || strings.Contains(string(p), "#") {
		return errors.New("illegal Charater(s)? Must not contains (){}[]#?")
	}
	return
}

func (p Password) Compare(p2 Password) (err error) {
	if ok := p != p2; ok {
		return errors.New("password is not correct?")
	}
	return nil
}

// NewUser() returns a pointer to a user with the given parameters.
func NewUser(firstname, lastname, nation Name, gendre Gendre, gmail Gmail, password Password, day, month, year int, id uuid.UUID) *User {
	if gendre == "female" {
		return &User{
			Firstname: firstname,
			Lastname:  lastname,
			Nation:    nation,
			Birthday: Birthday{
				Day:   Day(day),
				Month: Month(month),
				Year:  Year(year),
			},
			Gendre:   gendre,
			Gmail:    gmail,
			Password: password,
			VIPSubscription: VIPSubscription{
				lasts:        0,
				subscription: false,
			},
			ID:           id,
			PhotoProfile: defaultPPPath + "female.png",
		}
	} else {
		return &User{
			Firstname: firstname,
			Lastname:  lastname,
			Nation:    nation,
			Birthday: Birthday{
				Day:   Day(day),
				Month: Month(month),
				Year:  Year(year),
			},
			Gendre:   gendre,
			Gmail:    gmail,
			Password: password,
			VIPSubscription: VIPSubscription{
				lasts:        0,
				subscription: false,
			},
			ID:           id,
			PhotoProfile: defaultPPPath + "male.png",
		}
	}

}

func Sbirthday(birthday string) (Birthday, error) {
	slice := strings.Split(birthday, "-")
	if len(slice) != 3 {
		return Birthday{}, errors.New("not enough or many birthday information?")
	}
	day, err := strconv.Atoi(slice[2])
	if err != nil {
		return Birthday{}, err
	}
	month, err := strconv.Atoi(slice[1])
	if err != nil {
		return Birthday{}, err
	}
	year, err := strconv.Atoi(slice[0])
	if err != nil {
		return Birthday{}, err
	}

	return Birthday{
		Day:   Day(day),
		Month: Month(month),
		Year:  Year(year),
	}, nil
}
