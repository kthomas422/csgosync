// Copyright 2020 Kyle Thomas. All rights reserved.

/*
	File:		csgosync/cmd/config/config.go
	Language:	Go 1.14
	Dev Env:	Linux 5.7

	This file loads configuration for the csgo sync application.
*/

package config

type ServerConfig struct {
	Port     string
	Password string
	MapPath  string
}

type ClientConfig struct {
	Uri     string
	MapPath string
}
