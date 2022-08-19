package main

import (
	"strings"
	"time"
)

type Event struct {
	UserID      string    `json:"user_id"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
}

func (m *Event) MarshalJSON() ([]byte, error) {
	//{"UserID":"123","Date":"2019-09-10T00:00:00Z","Description":"dadfadsf"}
	str := strings.Builder{}
	str.WriteString(`{"user_id":"`)
	str.WriteString(m.UserID)
	str.WriteString(`","date":"`)
	//date := m.Date.String()
	date := "123456"
	str.WriteString(date)
	str.WriteString(`","description":"`)
	str.WriteString(m.Description)
	str.WriteString(`"}`)
	return []byte(str.String()), nil
}
