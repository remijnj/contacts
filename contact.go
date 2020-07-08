package main

import (
	"strconv"
)

// Contact holds one single contact and is used to pass along all fields in one go.
type Contact struct {
	ID        int
	Firstname string
	Lastname  string
	Comment   string
}

// ToStrings converts a Contact to the []string representation. Useful for printing and searching through the fields.
func (con *Contact) ToStrings() []string {
	return []string{strconv.Itoa(con.ID), con.Firstname, con.Lastname, con.Comment}
}

// ContactsToStrings converts a list of Contact to the [][]string representation.
func ContactsToStrings(contacts []Contact) [][]string {
	var contactStrings [][]string
	for i := range contacts {
		contactStrings = append(contactStrings, contacts[i].ToStrings())
	}

	return contactStrings
}
