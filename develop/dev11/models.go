package main

import (
	"strconv"
	"strings"
	"time"
)

type Event struct {
	UserID      string    `json:"user_id"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
}

func (m *Event) MarshalJSON() ([]byte, error) {
	//{"UserID":"123","Date":"2006-01-02","Description":"dadfadsf"}
	str := strings.Builder{}
	str.WriteString(`{"user_id":"`)
	str.WriteString(m.UserID)
	str.WriteString(`","date":"`)
	year, month, day := m.Date.Date()
	str.WriteString(strconv.Itoa(year))
	str.WriteRune('-')
	str.WriteString(strconv.Itoa(int(month)))
	str.WriteRune('-')
	str.WriteString(strconv.Itoa(day))
	str.WriteString(`","description":"`)
	str.WriteString(m.Description)
	str.WriteString(`"}`)
	return []byte(str.String()), nil
}
