package main

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

//маршалим дату по известным году, месяцу и дню
//в формате "2006-01-02"
func dateBuilder(str *strings.Builder, year, month, day int) {
	str.WriteString(strconv.Itoa(year))
	str.WriteRune('-')
	if month < 10 {
		str.WriteRune('0')
	}
	str.WriteString(strconv.Itoa(month))
	str.WriteRune('-')
	if day < 10 {
		str.WriteRune('0')
	}
	str.WriteString(strconv.Itoa(day))
}

//MarshalJSON - реализация интерфейса для события
func (m *Event) MarshalJSON() ([]byte, error) {
	//{"user_id":"123","date":"2019-09-10", "event_num":"1", "description":"dadfadsf"}
	str := strings.Builder{}
	str.WriteString(`{"user_id":"`)
	str.WriteString(m.UserID)
	str.WriteString(`" ,"date":"`)
	year, month, day := m.Date.Date()
	dateBuilder(&str, year, int(month), day)
	str.WriteString(`","event_num":"`)
	str.WriteString(strconv.Itoa(m.Num))
	str.WriteString(`","description":"`)
	str.WriteString(m.Description)
	str.WriteString(`"}`)
	return []byte(str.String()), nil
}

//UnmarshalJSON - реализация интерфейса для события
func (m *Event) UnmarshalJSON(json []byte) error {
	//{"user_id":"123","date":"2019-09-10","event_num":"1","description":"dadfadsf"}
	strJSON := string(json)
	trimmed := strings.Trim(strJSON, "{}\"")
	splitJSON := strings.Split(trimmed, "\",\"")
	keyVal := make([][]string, len(splitJSON))
	for i, substr := range splitJSON {
		keyVal[i] = strings.Split(substr, "\":\"")
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

	if strings.Contains(keyVal[2][0], "event_num") {
		var err error
		m.Num, err = strconv.Atoi(keyVal[2][1])
		if err != nil {
			return err
		}
	} else {
		return errors.New("error in event_num")
	}

	if strings.Contains(keyVal[3][0], "description") {
		m.Description = keyVal[3][1]
	} else {
		return errors.New("error in description")
	}

	return nil
}

//MarshalJSON - реализация интерфейса для ответа-ошибки
func (e ErrorResponse) MarshalJSON() ([]byte, error) {
	str := strings.Builder{}
	str.WriteString(`{"error":"`)
	str.WriteString(e.ErrorMsg)
	str.WriteString(`"}`)
	return []byte(str.String()), nil
}

//MarshalJSON - реализация интерфейса для ответа-результата
func (r ResultResponse) MarshalJSON() ([]byte, error) {
	var result strings.Builder
	result.WriteString(`{"result":"`)
	body, err := r.Result.MarshalJSON()
	if err != nil {
		return []byte{}, err
	}
	result.Write(body)
	result.WriteString(`"}`)
	return []byte(result.String()), nil
}

//MarshalJSON - реализация интерфейса для структуры,
//содержащей результат гет-запроса
func (g GetResult) MarshalJSON() ([]byte, error) {
	var result strings.Builder
	result.WriteString(`{"events":[`)
	for i, event := range g.Events {
		body, err := event.MarshalJSON()
		if err != nil {
			return []byte{}, err
		}
		result.Write(body)
		if i != len(g.Events)-1 {
			result.WriteRune(',')
		}
	}
	result.WriteString(`]}`)
	return []byte(result.String()), nil
}

//MarshalJSON - реализация интерфейса для структуры,
//содержащей результат пост-запроса
func (c CRUDResult) MarshalJSON() ([]byte, error) {
	var result strings.Builder
	result.WriteString(`{"event_num":`)
	result.WriteString(strconv.Itoa(c.Num))
	result.WriteString(`,"success":`)
	if c.Success {
		result.WriteString(`true`)
	} else {
		result.WriteString(`false`)
	}
	result.WriteRune('}')
	return []byte(result.String()), nil
}
