// Copyright 2020 Kyle Thomas. All rights reserved.

/*
	File:		csgosync/internal/auth/auth.go
	Language:	Go 1.14
	Dev Env:	Linux 5.7

	This file contains the functions for authenticating with the csgo sync application.
*/

package auth

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var creds struct {
	password string
	uri      string
}

func GetUri() error {
	uri, err := getInput("please enter the uri:")
	creds.uri = uri
	return err
}

func GetPass() error {
	pass, err := getInput("please enter the password:")
	creds.password = pass
	return err
}

func Wait() {
	_, err := getInput("done, press any key to continue:")
	if err != nil {
		fmt.Println("error: ", err)
	}
}

func getInput(prompt string) (string, error) {
	var (
		input string
		err   error
	)
	reader := bufio.NewReader(os.Stdin)

	fmt.Print(prompt + " ")
	input, err = reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	input = sanitizeInput(input)
	return input, nil
}

func Password() string {
	return creds.password
}

func Uri() string {
	return creds.uri
}

// sanitizeInput strips newline and carriage return from user input
func sanitizeInput(input string) string {
	input = strings.Replace(input, "\r", "", -1)
	input = strings.Replace(input, "\n", "", -1)
	return input
}
