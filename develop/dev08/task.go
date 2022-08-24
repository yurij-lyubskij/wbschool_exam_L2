package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	pr "github.com/mitchellh/go-ps"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

/*
=== Взаимодействие с ОС ===

Необходимо реализовать собственный шелл

встроенные команды: cd/pwd/echo/kill/ps
поддержать fork/exec команды
конвейер на пайпах

*/

//print working directory
func pwd(w io.Writer) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, wd)
	return nil
}

//change directory
func cd(args []string) error {
	if len(args) < 1 {
		return errors.New("no directory mentioned")
	}
	str := args[0]
	err := os.Chdir(str)
	if err != nil {
		return err
	}
	return nil
}

//echo arguments to output
func echo(args []string, w io.Writer) error {
	for _, str := range args {
		_, err := fmt.Fprintln(w, str)
		if err != nil {
			return err
		}
	}
	return nil
}

//kill process by pid
func kill(args []string) error {
	//проверяем, что есть аргументы
	if len(args) < 1 {
		return errors.New("no pid found")
	}
	//получаем PID из строки
	pid, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}
	//находим процесс
	p, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	//убиваем процесс
	err = p.Kill()
	if err != nil {
		return err
	}
	return nil
}

//выполняем внешнюю команду,
//привязывая райтер к выводу и ридер к вводу
func execute(args []string, w io.Writer, r io.Reader) error {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = r
	cmd.Stdout = w
	cmd.Stderr = w
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

//выводим процессы в райтер
//по умолчанию, только процессы,
//порожденные нашей программой
func ps(w io.Writer) error {
	all := false
	//получаем пид текущего процесса
	mypid := os.Getpid()
	//получаем список процессов
	processList, err := pr.Processes()
	if err != nil {
		return err
	}
	//выводим список процессов
	for _, p := range processList {
		if p.Pid() == mypid || p.PPid() == mypid || all {
			fmt.Fprintf(w, "PID= %d, PPID=  %d, CMD= %s\n", p.Pid(), p.PPid(), p.Executable())
		}
	}
	return nil
}

//выходим
func quit() {
	os.Exit(0)
}

//выполняем команду, передавая аргументы
func runCommand(args []string, w io.Writer, r io.Reader) error {
	if len(args) < 1 {
		return errors.New("no command")
	}
	name := args[0]
	args = args[1:]
	//проверяем соответствие названия команды
	//и вызываем нужную функцию
	switch name {
	case "cd":
		return cd(args)
	case "pwd":
		return pwd(w)
	case "echo":
		return echo(args, w)
	case "ps":
		return ps(w)
	case "kill":
		return kill(args)
	case "exec":
		return execute(args, w, r)
	case "\\quit":
		quit()
	default:
		return errors.New("unknown command")
	}
	return nil
}

//наша shell-утилита
//на входе - файлы для чтения и записи
func shell(f1 io.Reader, f2 io.Writer) {
	//буфер-writer
	var buf bytes.Buffer
	//создаем сканер для чтения из файла
	sc := bufio.NewScanner(f1)
	//читаем построчно
	for sc.Scan() {
		line := sc.Text()
		//разделяем пайп на команды
		pipe := strings.Split(line, "|")
		//выполняем команды пайпа
		for _, command := range pipe {
			//записываем вывод предыдущей команды в строку
			out := buf.String()
			buf.Reset()
			//разделяем аргументы
			args := strings.Split(command, " ")
			//создаем ридер из вывода прошлой команды
			reader := strings.NewReader(out)
			//выполянем новую команду
			err := runCommand(args, &buf, reader)
			if err != nil {
				fmt.Fprintln(f2, err)
			}
		}
		//выводим результат последней команды в пайпе
		fmt.Fprintf(f2, buf.String())
		buf.Reset()
	}
}

func main() {
	shell(os.Stdin, os.Stdout)
}
