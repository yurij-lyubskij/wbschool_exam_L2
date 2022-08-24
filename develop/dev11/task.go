package main

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
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
var storage Repository
var logger *log.Logger

//подготавливаем конфиг к чтению,
//репозиторий к сохранению данных
//и готовимся записывать логи в файл
func init() {
	storage = NewStorage()
	err := ReadConfig()
	if err != nil {
		log.Fatal(err.Error())
	}
	logger = log.Default()
	logfile := viper.GetString("log")
	file, err := os.Create(logfile)
	if err != nil {
		log.Fatal(err.Error())
	}
	logger.SetOutput(file)
}

//ReadConfig - Функция для подготовки
//чтения из конфига
func ReadConfig() error {
	viper.SetConfigName("config.json")
	viper.AddConfigPath("./")
	viper.SetConfigType("json")

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	return nil
}

func main() {
	//роутинг с помощью http.NewServeMux()
	//передаем url pattern и хэндлер
	mux := http.NewServeMux()
	mux.HandleFunc("/create_event", createHandler)
	mux.HandleFunc("/update_event", updateHandler)
	mux.HandleFunc("/delete_event", deleteHandler)
	mux.HandleFunc("/events_for_day", dayEventsHandler)
	mux.HandleFunc("/events_for_week", weekEventsHandler)
	mux.HandleFunc("/events_for_month", monthEventsHandler)
	mux.HandleFunc("/", defaultHandler)
	//используем Middleware
	siteHandler := logMiddleware(mux)
	//считываем порт из конфига
	address := viper.GetString("port")
	fmt.Println("launching server at port", address)
	//запускаем сервер
	http.ListenAndServe(address, siteHandler)

}
