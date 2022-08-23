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

func findLinks(n *html.Node, s *strings.Builder) {
	if n.Type == html.ElementNode {
		for _, a := range n.Attr {
			if a.Key == "href" {
				fmt.Fprintln(s, a.Val)
				break
			}
		}
	}
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

var maxLevel int
var dir string

func init() {
	dir, _ = os.Getwd()
}

func recursiveDownload(link string, level int, name string) {
	baseLink, err := url.Parse(link)
	if err != nil {
		log.Println("parse ", link, err.Error())
		return
	}
	resp, err := http.Get(link)
	if err != nil {
		log.Println("get ", link, err.Error())
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	fmt.Println(name)
	file, err := os.Create(name)
	if err != nil {
		log.Println("create, name = ", name, "err= ", err.Error(), " link =", link)
		return
	}
	defer file.Close()

	if level == maxLevel {
		file.Write(body)
		return
	}
	reader := bytes.NewReader(body)
	doc, err := html.Parse(reader)
	if err != nil {
		log.Println("htmlparse", err)
		return
	}
	s := strings.Builder{}
	findLinks(doc, &s)
	links := strings.Split(s.String(), "\n")
	for _, oneLink := range links {
		if oneLink == "/" {
			continue
		}
		parsedLink, err := url.Parse(oneLink)
		if err != nil {
			log.Println("parseOne ", oneLink, err.Error())
			return
		}
		extension := filepath.Ext(parsedLink.Path)
		if extension == "" {
			extension = ".html"
		}
		nextName := uuid.NewString() + extension
		nextURL := baseLink.ResolveReference(parsedLink)
		replaceLink := "file://" + dir + "/" + nextName
		fmt.Println(nextURL.String())
		body = bytes.Replace(body, []byte(oneLink), []byte(replaceLink), 1)
		recursiveDownload(nextURL.String(), level+1, nextName)
	}
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
	link := args[0]
	fname := "index.html"
	recursiveDownload(link, 1, fname)
}
