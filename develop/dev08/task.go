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

func pwd(w io.Writer) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	fmt.Fprintln(w, wd)
	return nil
}

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

func echo(args []string, w io.Writer) error {
	_, err := fmt.Fprintln(w, args)
	return err
}

func kill(args []string) error {
	if len(args) < 1 {
		return errors.New("no pid found")
	}
	pid, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}
	p, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	err = p.Kill()
	if err != nil {
		return err
	}
	return nil
}

func execute(args []string, w io.Writer, r io.Reader) error {
	cmd := exec.Command(args[0], args[1:]...)
	fmt.Println(args)
	cmd.Stdin = r
	cmd.Stdout = w
	cmd.Stderr = w
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func ps(w io.Writer) error {
	all := false
	mypid := os.Getpid()
	processList, err := pr.Processes()
	if err != nil {
		return err
	}

	for _, p := range processList {
		if p.Pid() == mypid || p.PPid() == mypid || all {
			fmt.Fprintf(w, "PID= %d, PPID=  %d, CMD= %s\n", p.Pid(), p.PPid(), p.Executable())
		}
	}
	return nil
}

func quit() {
	os.Exit(0)
}

func runCommand(args []string, w io.Writer, r io.Reader) error {
	name := args[0]
	args = args[1:]
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

func shell(f1 *os.File, f2 *os.File) {
	var buf bytes.Buffer
	sc := bufio.NewScanner(f1)
	for sc.Scan() {
		line := sc.Text()
		pipe := strings.Split(line, "|")
		for _, command := range pipe {
			out := buf.String()
			buf.Reset()
			args := strings.Split(command, " ")
			reader := strings.NewReader(out)
			err := runCommand(args, &buf, reader)
			if err != nil {
				fmt.Fprintln(f2, err)
			}
		}
		fmt.Fprintln(f2, buf.String())
		buf.Reset()
	}
}

func main() {
	shell(os.Stdin, os.Stdout)
}
