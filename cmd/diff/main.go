package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	instadiff "github.com/xabi93/instagram-diff"

	"github.com/xabi93/instagram-diff/instagram"
	"github.com/xabi93/instagram-diff/server"
)

func main() {
	var a App

	a.Init()

	a.Run()
}

type Conf struct {
	user        string
	password    string
	sessionFile string
	port        string
}

type App struct {
	cfg Conf
}

func (a *App) Init() {
	flag.StringVar(&a.cfg.user, "user", "", "Instagram username")
	flag.StringVar(&a.cfg.password, "password", "", "Instagram password")

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	flag.StringVar(&a.cfg.sessionFile, "sessionFile", fmt.Sprintf("%s/.instadiff", home), "Insta diff session file")

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

	fmt.Println("Login...")

	i, err = instagram.Login(a.cfg.user, a.cfg.password, a.cfg.sessionFile)
	if err != nil {
		return nil, err
	}

	return instadiff.New(i), nil
}

func (a App) restore() (*instagram.Instagram, error) {
	if _, err := os.Stat(a.cfg.sessionFile); os.IsNotExist(err) {
		fmt.Println("Session file does not exist, skipping...")
		return nil, nil
	}

	fmt.Printf("Restoring session from %s\n", a.cfg.sessionFile)

	return instagram.RestoreSession(a.cfg.sessionFile)
}
