// Copyright 2020 Kyle Thomas. All rights reserved.

/*
	File:		csgosync/cmd/config/config.go
	Language:	Go 1.15
	Dev Env:	Linux 5.9

	This file loads configuration for the csgo sync application.
*/

package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Both the client and server will have these values in their config
type baseConfig struct {
	Pass    string // Password for accessing the api
	MapPath string // Path to where the maps are stored
}

// Server configuration values
type ServerConfig struct {
	Port    string // Port to listen on
	LogFile string // Where to put logs
	*baseConfig
}

// Client configuration values
type ClientConfig struct {
	Uri string // Where the server is located
	*baseConfig
}

// Returns a populated baseConfig structure
func initConfig() *baseConfig {
	return &baseConfig{
		Pass:    viper.GetString("PASSWORD"),
		MapPath: viper.GetString("MAP_PATH"),
	}
}

// Returns a populated ServerConfig structure
func InitServerConfig() *ServerConfig {
	return &ServerConfig{
		viper.GetString("PORT"),
		viper.GetString("LOG_FILE"),
		initConfig(),
	}
}

// Returns a populated ClientConfig structure
func InitClientConfig() *ClientConfig {
	c := &ClientConfig{
		viper.GetString("URI"),
		initConfig(),
	}
	if !strings.HasPrefix(c.Uri, "http://") {
		c.Uri = "http://" + c.Uri
	}
	return c
}

// Prompts the user to enter the URI
func (c *ClientConfig) GetUri() error {
	uri, err := getInput("please enter the uri:")
	if !strings.HasPrefix(uri, "http://") {
		uri = "http://" + uri
	}
	c.Uri = uri
	return err
}

// Prompts the user to enter the password
func (c *baseConfig) GetPass() error {
	pass, err := getInput("please enter the password:")
	c.Pass = pass
	return err
}

// Since winturds closes the cmd when it exits "wait" for user so they can read the output.
func Wait() {
	_, err := getInput("done, press \"enter\" to continue")
	if err != nil {
		fmt.Println("error: ", err)
	}
}

// Displays the prompt to the user and gather's their response
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

// sanitizeInput strips newline and carriage return from user input
func sanitizeInput(input string) string {
	input = strings.Replace(input, "\r", "", -1)
	input = strings.Replace(input, "\n", "", -1)
	return input
}
