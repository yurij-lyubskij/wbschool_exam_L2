package main

import (
	"errors"
	"strings"
	"time"
)

//В чистой архитектуре принято разделение на
//delivery, usecase и repository
//Но в данной задаче вся бизнес-логика
//сводится к сохранению и изменению событий
//Поэтому выделяю только delivery для хэндлеров
//и repository для хранения. Для данной задачи
//бизнес-логика в repository

//day []descriptions
type daymap map[int][]string

//month dayMap
type monthmap map[int]daymap

//year monthMap
type yearMap map[int]monthmap

//user_id yearMap
type EventsMap map[string]yearMap

type Storage struct {
	Events EventsMap
}

func NewStorage() Repository {
	return &Storage{
		Events: make(EventsMap),
	}
}

func (m *Storage) createEvent(event Event) (num int, err error) {
	id := event.UserID
	if m.Events[id] == nil {
		m.Events[id] = make(yearMap)
	}
	year, month, day := event.Date.Date()
	if m.Events[id][year] == nil {
		m.Events[id][year] = make(monthmap)
	}
	if m.Events[id][year][int(month)] == nil {
		m.Events[id][year][int(month)] = make(daymap)
	}
	arr := m.Events[id][year][int(month)][day]
	arr = append(arr, event.Description)
	m.Events[id][year][int(month)][day] = arr
	num = len(arr)
	return
}

func (m *Storage) updateEvent(event Event) error {
	id := event.UserID
	err := errors.New("event does not exist")
	if m.Events[id] == nil {
		return err
	}
	year, month, day := event.Date.Date()
	if m.Events[id][year] == nil {
		return err
	}
	if m.Events[id][year][int(month)] == nil {
		return err
	}
	arr := m.Events[id][year][int(month)][day]
	if len(arr) < event.Num {
		return err
	}
	arr[event.Num-1] = event.Description
	return nil
}

func (m *Storage) deleteEvent(event Event) error {
	id := event.UserID
	err := errors.New("event does not exist")
	if m.Events[id] == nil {
		return err
	}
	year, month, day := event.Date.Date()
	if m.Events[id][year] == nil {
		return err
	}
	if m.Events[id][year][int(month)] == nil {
		return err
	}
	arr := m.Events[id][year][int(month)][day]
	if len(arr) < event.Num {
		return err
	}
	arr[event.Num-1] = event.Description
	arr = append(arr[:event.Num-1], arr[event.Num:]...)
	m.Events[id][year][int(month)][day] = arr
	return nil
}

func (m *Storage) getEventsByDay(id string, date time.Time) []Event {
	year, month, day := date.Date()
	var result []Event
	if m.Events[id] == nil {
		return result
	}
	if m.Events[id][year] == nil {
		return result
	}
	if m.Events[id][year][int(month)] == nil {
		return result
	}
	arr := m.Events[id][year][int(month)][day]
	event := Event{
		UserID:      id,
		Date:        date,
		Num:         0,
		Description: "",
	}

	for i, description := range arr {
		event.Num = i + 1
		event.Description = description
		result = append(result, event)
	}

	return result
}

func (m *Storage) getEventsByMonth(id string, date time.Time) []Event {
	year, month, _ := date.Date()
	var result []Event
	if m.Events[id] == nil {
		return result
	}
	if m.Events[id][year] == nil {
		return result
	}
	monthInMap := m.Events[id][year][int(month)]
	event := Event{
		UserID:      id,
		Date:        date,
		Num:         0,
		Description: "",
	}
	for _, arr := range monthInMap {
		for i, description := range arr {
			event.Num = i + 1
			event.Description = description
			result = append(result, event)
		}
	}
	return result
}

func (m *Storage) getEventsByWeek(id string, date time.Time) []Event {
	year, month, _ := date.Date()
	_, w1 := date.ISOWeek()
	var result []Event
	if m.Events[id] == nil {
		return result
	}
	if m.Events[id][year] == nil {
		return result
	}
	monthInMap := m.Events[id][year][int(month)]
	event := Event{
		UserID:      id,
		Date:        date,
		Num:         0,
		Description: "",
	}
	str := strings.Builder{}
	for day, arr := range monthInMap {
		dateBuilder(&str, year, int(month), day)
		t, _ := time.Parse(dateForm, str.String())
		str.Reset()
		_, w2 := t.ISOWeek()
		if w1 == w2 {
			for i, description := range arr {
				event.Num = i + 1
				event.Description = description
				result = append(result, event)
			}
		}
	}
	return result
}
