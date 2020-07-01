package auth

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var globals struct {
	password string
	uri      string
}

func GetUserCreds() error {
	var err error
	fmt.Print("Please enter the server url: ")
	reader := bufio.NewReader(os.Stdin)
	globals.uri, err = reader.ReadString('\n')
	if err != nil {
		return err
	}
	globals.uri = strings.Replace(globals.uri, "\r\n", "", -1)

	fmt.Print("Please enter the password: ")
	globals.password, err = reader.ReadString('\n')
	if err != nil {
		return err
	}
	globals.password = strings.Replace(globals.password, "\r\n", "", -1)
	return nil
}

func Password() string {
	return globals.password
}

func Uri() string {
	return globals.uri
}
