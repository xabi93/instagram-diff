package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"syscall"

	instadiff "github.com/xabi93/instagram-diff"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/xabi93/instagram-diff/instagram"
	"github.com/xabi93/instagram-diff/server"
)

func main() {
	var a App

	a.Init()
	a.Run()
}

type Conf struct {
	sessionFile string
	port        string
}

type App struct {
	cfg Conf
}

func (a *App) Init() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	flag.StringVar(&a.cfg.sessionFile, "sessionFile", fmt.Sprintf("%s/.instadiff", home), "Session file")

	flag.StringVar(&a.cfg.port, "port", "3000", "Port to expose result")

	flag.Parse()
}

func (a App) Run() {
	i, err := a.login()
	if err != nil {
		log.Fatal(err)
	}

	if err := server.Serve(a.cfg.port, i); err != nil {
		log.Fatal(err)
	}
}

func (a App) login() (*instadiff.Instadiff, error) {
	i, err := a.restore()
	if err != nil {
		return nil, err
	}

	if i != nil {
		return instadiff.New(i), nil
	}

	user, pass, err := a.askUserPass()

	fmt.Println("Login...")

	i, err = instagram.Login(user, pass, a.cfg.sessionFile)
	if err != nil {
		return nil, err
	}

	return instadiff.New(i), nil
}

func (a App) restore() (*instagram.Instagram, error) {
	if _, err := os.Stat(a.cfg.sessionFile); os.IsNotExist(err) {
		return nil, nil
	}

	fmt.Printf("Restoring session from %s\n", a.cfg.sessionFile)

	i, err := instagram.RestoreSession(a.cfg.sessionFile)
	if err != nil {
		return nil, err
	}

	err = i.Ping()
	if err == nil {
		return i, nil
	}
	if errors.As(err, &instagram.AuthError{}) {
		fmt.Println("Session outdated")
		os.Remove(a.cfg.sessionFile)
		return nil, nil
	}

	return i, nil
}

func (App) askUserPass() (string, string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Username: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		return "", "", err
	}

	fmt.Print("Enter Password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", "", err
	}

	fmt.Print("\n")

	return username, string(bytePassword), nil
}
