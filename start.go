package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strconv"
)

func main() {
	// print file and line numbers in log lines
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	db := dbInit()
	defer db.Close()

	//dbAddContact(db)
	//dbAddTestContacts(db)
	//dbDebugQuery(db)
	headings, contacts, _ := dbContactList(db)

	ContactsApp(headings, contacts, func(con Contact) {
		dbSaveContact(db, con)
	})
}

//////////////////////////////////////////////
//
// DB stuff
//
func dbContactList(database *sql.DB) (headers []string, contacts []Contact, err error) {
	rows, err := database.Query("SELECT id, firstname, lastname, comment FROM people")
	if err != nil {
		log.Fatal(err)
	}
	var con Contact
	//var id int
	//var firstname string
	//var lastname string
	//var comment string
	var items []Contact

	headers = []string{"ID", "First name", "Last name", "Comment"}
	for rows.Next() {
		rows.Scan(&con.ID, &con.Firstname, &con.Lastname, &con.Comment)
		items = append(items, con)
	}

	return headers, items, err
}

func dbSaveContact(database *sql.DB, contact Contact) error {
	fmt.Println("dbSaveContact id=" + strconv.Itoa(contact.ID))
	if contact.ID > 0 {
		return dbUpdateContact(database, contact)
	}

	return dbAddContact(database, contact)
}

func dbUpdateContact(database *sql.DB, contact Contact) error {
	fmt.Println("updating " + contact.Firstname + " " + contact.Lastname + " (" + contact.Comment + ")")
	statement, err :=
		database.Prepare("UPDATE people SET firstname=?, lastname=?, comment=? WHERE id=?")
	if err != nil {
		log.Fatal(err)
	}
	_, err = statement.Exec(contact.Firstname, contact.Lastname, contact.Comment, contact.ID)
	if err != nil {
		log.Fatal(err)
	}

	return err
}

func dbAddContact(database *sql.DB, contact Contact) error {
	fmt.Println("adding " + contact.Firstname + " " + contact.Lastname + " (" + contact.Comment + ")")
	statement, err :=
		database.Prepare("INSERT INTO people (firstname, lastname, comment) VALUES (?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = statement.Exec(contact.Firstname, contact.Lastname, contact.Comment)
	if err != nil {
		log.Fatal(err)
	}

	return err
}

func dbInit() *sql.DB {
	database, err := sql.Open("sqlite3", "./contacts.db")
	if err != nil {
		log.Fatal(err)
	}

	// add the table
	statement, err :=
		database.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, firstname TEXT, lastname TEXT, comment TEXT)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = statement.Exec()
	if err != nil {
		log.Fatal(err)
	}

	// set the user_version
	statement, err = database.Prepare("PRAGMA user_version = 1")
	if err != nil {
		log.Fatal(err)
	}
	_, err = statement.Exec()
	if err != nil {
		log.Fatal(err)
	}

	dbMigrate(database)

	return database
}

func dbMigrate(database *sql.DB) {
	// get the user_version and if needed run DB migration
	row := database.QueryRow("PRAGMA user_version")
	var userVersion int
	row.Scan(&userVersion)

	switch userVersion {
	case 0:
		// add comment field
		statement, err := database.Prepare("ALTER TABLE people ADD comment TEXT")
		if err != nil {
			log.Fatal(err)
		}
		_, err = statement.Exec()
		if err != nil {
			log.Fatal(err)
		}

		// set user_version to 1
		statement, err = database.Prepare("PRAGMA user_version = 1")
		if err != nil {
			log.Fatal(err)
		}
		_, err = statement.Exec()
		if err != nil {
			log.Fatal(err)
		}
		fallthrough
	case 1:
		// current version
	}
}

////////////////////
//
// Debug / Old stuff
//
func dbDebugQuery(database *sql.DB) error {
	rows, err := database.Query("SELECT id, firstname, lastname, comment FROM people")
	if err != nil {
		log.Fatal(err)
	}
	var id int
	var firstname string
	var lastname string
	var comment string
	for rows.Next() {
		rows.Scan(&id, &firstname, &lastname, &comment)
		fmt.Println(strconv.Itoa(id) + ": " + firstname + " " + lastname + " " + comment)
	}

	return err
}
