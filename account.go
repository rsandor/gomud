package main

import (
	//"os"
	"fmt"
	//"encoding/xml"
)

type Account struct {
	XMLName		string		`xml:"user"`
	Name		[]byte		`xml:"name,attr"`
	Password	[]byte		`xml:"password,attr"`
	Admin		bool		`xml:"admin,attr"`
}

type AccountNotFoundError string
func (err AccountNotFoundError) Error() string {
	return fmt.Sprintf("No account with the name '%s' exists.", string(err))
}

type IncorrectPasswordError string
func (err IncorrectPasswordError) Error() string {
	return fmt.Sprintf("Incorrect password for account with name '%s'", string(err))
}


// Attempts to load an account with the given name and password.
func LoadUser(name, password string) (*Account, error) {
	return nil, nil
}


