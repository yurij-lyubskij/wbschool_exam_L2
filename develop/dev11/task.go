package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

/*
=== HTTP server ===

Реализовать HTTP сервер для работы с календарем. В рамках задания необходимо работать строго со стандартной HTTP библиотекой.
В рамках задания необходимо:
	1. Реализовать вспомогательные функции для сериализации объектов доменной области в JSON.
	2. Реализовать вспомогательные функции для парсинга и валидации параметров методов /create_event и /update_event.
	3. Реализовать HTTP обработчики для каждого из методов API, используя вспомогательные функции и объекты доменной области.
	4. Реализовать middleware для логирования запросов
Методы API: POST /create_event POST /update_event POST /delete_event GET /events_for_day GET /events_for_week GET /events_for_month
Параметры передаются в виде www-url-form-encoded (т.е. обычные user_id=3&date=2019-09-09).
В GET методах параметры передаются через queryString, в POST через тело запроса.
В результате каждого запроса должен возвращаться JSON документ содержащий либо {"result": "..."} в случае успешного выполнения метода,
либо {"error": "..."} в случае ошибки бизнес-логики.

В рамках задачи необходимо:
	1. Реализовать все методы.
	2. Бизнес логика НЕ должна зависеть от кода HTTP сервера.
	3. В случае ошибки бизнес-логики сервер должен возвращать HTTP 503. В случае ошибки входных данных (невалидный int например) сервер должен возвращать HTTP 400. В случае остальных ошибок сервер должен возвращать HTTP 500. Web-сервер должен запускаться на порту указанном в конфиге и выводить в лог каждый обработанный запрос.
	4. Код должен проходить проверки go vet и golint.
*/

func main() {

	http.HandleFunc("/create_event", createHandler)
	http.HandleFunc("/update_event", updateHandler)
	http.HandleFunc("/delete_event", deleteHandler)
	http.HandleFunc("/events_for_day", dayEventsHandler)
	http.HandleFunc("/events_for_week", weekEventsHandler)
	http.HandleFunc("/events_for_month", monthEventsHandler)

	fmt.Println("starting server at :8080")
	t, err := time.Parse(dateForm, "2019-09-10")
	//d := time.Date(2000, 2, 1, 12, 30, 0, 0, time.UTC)
	year, month, day := t.Date()
	fmt.Println(year, month, day)
	if err != nil {
		fmt.Println(err)
	}
	event := Event{
		UserID:      "123",
		Date:        t,
		Num:         100,
		Description: "dadfadsf",
	}
	fmt.Println(event)
	bytes, err := json.Marshal(&event)
	fmt.Println(string(bytes))
	if err != nil {
		fmt.Println(err)
	}
	var newevent Event
	err = json.Unmarshal(bytes, &newevent)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(newevent)
	fmt.Println(time.Now())
	//http.ListenAndServe(":8080", nil)
	storage := NewStorage()
	_, err = storage.createEvent(event)
	if err != nil {
		fmt.Println(err)
	}
	event.Description = "new1"
	_, err = storage.createEvent(event)
	if err != nil {
		fmt.Println(err)
	}
	event.Description = "new2"
	_, err = storage.createEvent(event)
	if err != nil {
		fmt.Println(err)
	}
	event.Description = "new3"
	_, err = storage.createEvent(event)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(*storage)
	events := storage.getEventsByWeek(event.UserID, event.Date)
	fmt.Println(events)
}
