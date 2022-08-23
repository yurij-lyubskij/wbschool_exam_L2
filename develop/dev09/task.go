package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

/*
=== Утилита wget ===

Реализовать утилиту wget с возможностью скачивать сайты целиком

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

//находим ссылки в документе
func findLinks(n *html.Node, s *strings.Builder) {
	if n.Type == html.ElementNode {
		for _, a := range n.Attr {
			//если ссылка, пишем в буфер
			if a.Key == "href" {
				fmt.Fprintln(s, a.Val)
				break
			}
		}
	}
	//обходим узлы-детей рекурсивно
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		findLinks(c, s)
	}
}

var usageStr = `
Usage: wget [options] [URL with scheme]
file - stdin по умолчанию
Options:
-r - "recursive" - скачивать рекурсивно
`

func usage() {
	fmt.Printf("%s\n", usageStr)
	os.Exit(0)
}

//глубина рекурсии
var maxLevel int

//текущая директория
var dir string

func init() {
	//определяется текущая директория
	dir, _ = os.Getwd()
}

func recursiveDownload(link string, level int, name string) {
	//парсим ссылку, чтобы потом использовать ResolveReference
	baseLink, err := url.Parse(link)
	if err != nil {
		log.Println("parse ", link, err.Error())
		return
	}
	//выполняем GET запрос
	resp, err := http.Get(link)
	if err != nil {
		log.Println("get ", link, err.Error())
		return
	}
	//закрываем тело в конце
	defer resp.Body.Close()
	//вычитываем тело в слайс байт
	body, err := io.ReadAll(resp.Body)
	fmt.Println(name)
	//создаем файл с названием, переданным снаружи
	file, err := os.Create(name)
	if err != nil {
		log.Println("create, name = ", name, "err= ", err.Error(), " link =", link)
		return
	}
	//в конце закрываем файл
	defer file.Close()

	//если мы на максимальной глубине,
	//то ничего не подгружаем,
	//сохраняем страницу в файл
	if level == maxLevel {
		file.Write(body)
		return
	}
	//парсим страницу
	reader := bytes.NewReader(body)
	doc, err := html.Parse(reader)
	if err != nil {
		log.Println("htmlparse", err)
		return
	}
	//находим ссылки и сохраняем в строку
	s := strings.Builder{}
	findLinks(doc, &s)
	//разделяем ссылки
	links := strings.Split(s.String(), "\n")
	//заменяем каждую ссылку в исходном файле
	//и загружаем файл рекрсивным вызовом функции
	for _, oneLink := range links {
		//не загружаем по второму разу
		//главную страницу
		if oneLink == "/" {
			continue
		}
		//парсим ссылку
		parsedLink, err := url.Parse(oneLink)
		if err != nil {
			log.Println("parseOne ", oneLink, err.Error())
			return
		}
		//получаем расширение для файла,
		//который хотим скачать
		extension := filepath.Ext(parsedLink.Path)
		if extension == "" {
			extension = ".html"
		}
		//файл сохраняем под случайным
		//уникальным именем,но с расширением
		nextName := uuid.NewString() + extension
		//получаем абсолютную ссылку,
		//даже если была относительная
		nextURL := baseLink.ResolveReference(parsedLink)
		//готовим локальную ссылку на замену
		replaceLink := "file://" + dir + "/" + nextName
		fmt.Println(nextURL.String())
		//заменяем ссылку в исходном файле
		body = bytes.Replace(body, []byte(oneLink), []byte(replaceLink), 1)
		//рекурсивно загружаем файл по ссылке и сохраняем
		recursiveDownload(nextURL.String(), level+1, nextName)
	}
	//после всех изменений ссылок сохраняем файл
	file.Write(body)
}

func main() {
	//настройки утилиты - возможные флаги
	flag.IntVar(&maxLevel, "r", 1, "загружать рекурсивно")
	flag.Usage = usage
	//парсим флаги
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		usage()
		return
	}
	//ссылка
	link := args[0]
	//главная страница
	fname := "index.html"
	//загружаем рекурсивно, начиная с 1 уровня глубины
	recursiveDownload(link, 1, fname)
}
