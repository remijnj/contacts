package main

import (
	"strconv"
)

type Contact struct {
	Id        int
	Firstname string
	Lastname  string
	Comment   string
}

func (con *Contact) ToStrings() []string {
	return []string{strconv.Itoa(con.Id), con.Firstname, con.Lastname, con.Comment}
}

func ContactsToStrings(contacts []Contact) [][]string {
	var contactStrings [][]string
	for i := range contacts {
		contactStrings = append(contactStrings, contacts[i].ToStrings())
	}

	return contactStrings
}
