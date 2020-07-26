// Copyright 2020 Kyle Thomas. All rights reserved.

/*
	File:		csgosync/internal/constants/constants.go
	Language:	Go 1.14
	Dev Env:	Linux 5.7

	This file contains constants for the csgo sync application (todo: move to config file)
*/

package constants

const (
	//ClientMapDir = "C:\\ProgramFilesx86\\Steam\\steamapps\\common\\csgo\\maps"
	ClientMapDir = "/home/kyle/steamapps2"
	ServerMapDir = "/home/kyle/steamapps" // users will have to make symlink to map folder
)
