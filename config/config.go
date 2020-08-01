// Copyright 2020 Kyle Thomas. All rights reserved.

/*
	File:		csgosync/cmd/config/config.go
	Language:	Go 1.14
	Dev Env:	Linux 5.7

	This file loads configuration for the csgo sync application.
*/

package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/kthomas422/csgosync/internal/logging"

	"github.com/spf13/viper"
)

type baseConfig struct {
	Pass    string
	MapPath string
}

type ServerConfig struct {
	Port    string
	LogFile string
	LogLvl  logging.Level
	*baseConfig
}

type ClientConfig struct {
	Uri string
	*baseConfig
}

func initConfig() *baseConfig {
	return &baseConfig{
		Pass:    viper.GetString("PASSWORD"),
		MapPath: viper.GetString("MAP_PATH"),
	}
}

func InitServerConfig() *ServerConfig {
	return &ServerConfig{
		viper.GetString("PORT"),
		viper.GetString("LOG_FILE"),
		logging.Level(viper.GetInt("LOG_LEVEL")),
		initConfig(),
	}
}

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

func (c *ClientConfig) GetUri() error {
	uri, err := getInput("please enter the uri:")
	if !strings.HasPrefix(uri, "http://") {
		uri = "http://" + uri
	}
	c.Uri = uri
	return err
}

func (c *baseConfig) GetPass() error {
	pass, err := getInput("please enter the password:")
	c.Pass = pass
	return err
}

func Wait() {
	_, err := getInput("done, press any key to continue")
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

// sanitizeInput strips newline and carriage return from user input
func sanitizeInput(input string) string {
	input = strings.Replace(input, "\r", "", -1)
	input = strings.Replace(input, "\n", "", -1)
	return input
}
