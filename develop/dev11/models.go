package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const dateForm = "2006-01-02"

type Event struct {
	UserID      string    `json:"user_id"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
}

func (m *Event) MarshalJSON() ([]byte, error) {
	//{"user_id":"123","date":"2019-09-10","description":"dadfadsf"}
	str := strings.Builder{}
	str.WriteString(`{"user_id":"`)
	str.WriteString(m.UserID)
	str.WriteString(`" ,"date":"`)
	year, month, day := m.Date.Date()
	str.WriteString(strconv.Itoa(year))
	str.WriteRune('-')
	if month < 10 {
		str.WriteRune('0')
	}
	str.WriteString(strconv.Itoa(int(month)))
	str.WriteRune('-')
	if day < 10 {
		str.WriteRune('0')
	}
	str.WriteString(strconv.Itoa(day))
	str.WriteString(`","description":"`)
	str.WriteString(m.Description)
	str.WriteString(`"}`)
	return []byte(str.String()), nil
}

func (m *Event) UnmarshalJSON(json []byte) error {
	//{"user_id":"123","date":"2019-09-10","description":"dadfadsf"}
	strJSON := string(json)
	trimmed := strings.Trim(strJSON, "{}\"")
	fmt.Println(trimmed)
	splitJSON := strings.Split(trimmed, "\",\"")
	fmt.Println(splitJSON)
	keyVal := make([][]string, len(splitJSON))
	for i, substr := range splitJSON {
		keyVal[i] = strings.Split(substr, "\":\"")
		fmt.Println(keyVal[i])
	}

	if strings.Contains(keyVal[0][0], "user_id") {
		m.UserID = keyVal[0][1]
	} else {
		return errors.New("error in user_id")
	}

	if strings.Contains(keyVal[1][0], "date") {
		t, err := time.Parse(dateForm, keyVal[1][1])
		if err != nil {
			return err
		}
		m.Date = t
	} else {
		return errors.New("error in date")
	}

	if strings.Contains(keyVal[2][0], "description") {
		m.Description = keyVal[2][1]
	} else {
		return errors.New("error in description")
	}

	return nil
}
