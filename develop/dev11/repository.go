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

//ключ - день, значение- []descriptions
type daymap map[int][]string

//ключ - месяц, значение- daymap
type monthmap map[int]daymap

//ключ - год, значение- monthmap
type yearMap map[int]monthmap

//EventsMap - карта, содержащая
//события. Список описаний
//событий дня получается,
//как map[user_id][год][мес][день]
type EventsMap map[string]yearMap

//Storage - Структура для хранения событий,
//для которой определены методы работы с ними
type Storage struct {
	Events EventsMap
}

//NewStorage - конструктор, создающий хранилище и
//возвращающий Repository. Repository - интерфейс
func NewStorage() Repository {
	return &Storage{
		Events: make(EventsMap),
	}
}

//создание события
func (m *Storage) createEvent(event Event) (num int, err error) {
	id := event.UserID
	//проверяем случаи, когда внутри nil
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
	//добавляем описание в массив и сохраняем
	arr := m.Events[id][year][int(month)][day]
	arr = append(arr, event.Description)
	m.Events[id][year][int(month)][day] = arr
	//возвращаем номер добавленного описания
	num = len(arr)
	return
}

//обновление события
func (m *Storage) updateEvent(event Event) error {
	id := event.UserID
	err := errors.New("event does not exist")
	//проверяем случаи, когда внутри nil
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
	//если нет такого по порядку события
	if len(arr) < event.Num {
		return err
	}
	//обновляем описание
	arr[event.Num-1] = event.Description
	return nil
}

//удаление события
func (m *Storage) deleteEvent(event Event) error {
	id := event.UserID
	err := errors.New("event does not exist")
	//проверяем случаи, когда внутри nil
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
	//удаляем событие
	arr = append(arr[:event.Num-1], arr[event.Num:]...)
	//сохраняем
	m.Events[id][year][int(month)][day] = arr
	return nil
}

//список событий за день
func (m *Storage) getEventsByDay(id string, date time.Time) []Event {
	year, month, day := date.Date()
	var result []Event
	//проверяем случаи, когда внутри nil
	if m.Events[id] == nil {
		return result
	}
	if m.Events[id][year] == nil {
		return result
	}
	if m.Events[id][year][int(month)] == nil {
		return result
	}
	//события за день
	arr := m.Events[id][year][int(month)][day]
	event := Event{
		UserID:      id,
		Date:        date,
		Num:         0,
		Description: "",
	}
	//составляем список событий по описаниям
	for i, description := range arr {
		event.Num = i + 1
		event.Description = description
		result = append(result, event)
	}

	return result
}

//список событий за месяц
func (m *Storage) getEventsByMonth(id string, date time.Time) []Event {
	year, month, _ := date.Date()
	var result []Event
	//проверяем случаи, когда внутри nil
	if m.Events[id] == nil {
		return result
	}
	if m.Events[id][year] == nil {
		return result
	}
	//события за месяц
	monthInMap := m.Events[id][year][int(month)]
	event := Event{
		UserID:      id,
		Date:        date,
		Num:         0,
		Description: "",
	}
	//составляем список событий по описаниям
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
	//сохраняем текущую неделю
	_, w1 := date.ISOWeek()
	var result []Event
	//проверяем случаи, когда внутри nil
	if m.Events[id] == nil {
		return result
	}
	if m.Events[id][year] == nil {
		return result
	}
	//события за месяц
	monthInMap := m.Events[id][year][int(month)]
	event := Event{
		UserID:      id,
		Date:        date,
		Num:         0,
		Description: "",
	}
	str := strings.Builder{}
	//составляем список событий по описаниям
	//обходим события месяца, проверяем,
	//произошли ли они на этой неделе
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
