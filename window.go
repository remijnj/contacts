package main

import (
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"strconv"
	"strings"
)

// ContactsApp starts the application and lists the contacts in a table with.
func ContactsApp(headers []Header, contacts []Contact, save func(Contact)) {
	app := app.New()
	window := app.NewWindow("Contacts")
	window.Resize(fyne.NewSize(400, 400))

	var table contactsTable

	var clickContact = func(con Contact) {
		fmt.Println("clickContact(" + strconv.Itoa(con.ID) + ")")
		editWindow := newEditWindow(app, &con, func(contact Contact) {
			// update it in the database
			save(contact)

			// now update it also in the UI
			table.update(contact)

		})
		editWindow.Show()
	}

	contactsBox := makeTable(headers, contacts, clickContact)
	contactsScroller := widget.NewVScrollContainer(contactsBox)

	table.Headers = headers
	table.Contacts = contacts
	table.Box = contactsBox
	table.Click = clickContact

	search := newSearchEntry(&table)
	search.SetPlaceHolder("Search")

	addbutton := widget.NewButton("Add", func() {
		addWindow := newEditWindow(app, nil, func(contact Contact) {
			// add it to the database
			save(contact)

			// now add it to the UI
			//  1. add in contacts
			//  2. re-query database and redo UI
			table.add(contact)

		})
		addWindow.Show()
	})

	content := fyne.NewContainerWithLayout(
		layout.NewBorderLayout(search, addbutton, nil, nil),
		search, addbutton, contactsScroller)
	window.SetContent(content)
	window.SetMaster()
	window.ShowAndRun()
}

func newEditWindow(app fyne.App, con *Contact, save func(Contact)) fyne.Window {
	window := app.NewWindow("Add")

	firstname := widget.NewEntry()
	lastname := widget.NewEntry()
	comment := widget.NewEntry()

	if con != nil {
		firstname.SetText(con.Firstname)
		lastname.SetText(con.Lastname)
		comment.SetText(con.Comment)
	} else {
		firstname.SetPlaceHolder("First name")
		lastname.SetPlaceHolder("Last name")
		comment.SetPlaceHolder("Comment")
	}

	savebutton := widget.NewButton("Save", func() {
		var contact Contact
		if (con != nil) {
			contact.ID = con.ID
		}
		contact.Firstname = firstname.Text
		contact.Lastname = lastname.Text
		contact.Comment = comment.Text
		save(contact)
		window.Close()
	})

	form := fyne.NewContainerWithLayout(
		layout.NewHBoxLayout(),
		firstname, lastname, comment, layout.NewSpacer(), savebutton)

	window.SetContent(form)

	return window
}

//////////////////////////////////////
//
// searchEntry
//
type searchEntry struct {
	widget.Entry
	search string
	table  *contactsTable
}

func (e *searchEntry) TypedKey(key *fyne.KeyEvent) {
	e.Entry.TypedKey(key)
	//fmt.Println("TypedKey:" + e.Text)
	if e.search != e.Text {
		e.search = e.Text
		e.table.filterContactsUI(e.Text)
	}
}

func (e *searchEntry) TypedRune(r rune) {
	e.Entry.TypedRune(r)
	//fmt.Println("TypedRune:" + e.Text)
	if e.search != e.Text {
		e.search = e.Text
		e.table.filterContactsUI(e.Text)
	}
}

func newSearchEntry(table *contactsTable) *searchEntry {
	entry := &searchEntry{}
	entry.table = table
	entry.ExtendBaseWidget(entry)
	return entry
}

//////////////////////////////////////
//
// contactsTable
//
type Header struct {
	text string
	display bool
}

type contactsTable struct {
	Headers      []Header
	Contacts      []Contact
	Box           *widget.Box
	CurrentFilter string // nice to keep current filter and re-filter when doing an add/update
	Click         func(con Contact)
}

func (c *contactsTable) update(contact Contact) {
	for i, con := range c.Contacts {
		if con.ID == contact.ID {
			c.Contacts[i] = contact
			break
		}
	}
	newContacts := filterContacts(c.Contacts, c.CurrentFilter)
	c.updateContacts(newContacts)
}

func (c *contactsTable) add(con Contact) {
	c.Contacts = append(c.Contacts, con)
	newContacts := filterContacts(c.Contacts, c.CurrentFilter)
	c.updateContacts(newContacts)
}

func (c *contactsTable) filterContactsUI(search string) {
	newContacts := filterContacts(c.Contacts, search)
	c.CurrentFilter = search
	c.updateContacts(newContacts)
}

func (c *contactsTable) updateContacts(newContacts []Contact) {
	newBox := makeTable(c.Headers, newContacts, c.Click)
	for i := 0; i < len(c.Box.Children); i++ {
		c.Box.Children[i] = newBox.Children[i]
	}
	c.Box.Refresh() // force re-draw (needed to make add() work)
}

func filterContacts(contacts []Contact, search string) (newContacts []Contact) {
	fmt.Println("filterContacts(" + strconv.Itoa(len(contacts)) + ", " + search + ")")
	// keep only contacts which match the search string in any field
	for i := 0; i < len(contacts); i++ {
		contactStrings := contacts[i].ToStrings()
		for j := 0; j < len(contactStrings); j++ {
			if strings.Contains(strings.ToLower(contactStrings[j]), strings.ToLower(search)) {
				newContacts = append(newContacts, contacts[i])

				// stop searching after you found a hit on this contact and added it
				// without this break it will add the contact once for each field hit
				break
			}
		}
	}

	fmt.Println("filterContacts return=" + strconv.Itoa(len(newContacts)) + " contacts")
	return newContacts
}

func headersToText(headers []Header, displayOnly bool) []string {
	var headerStrings []string
	for _, header := range(headers) {
		if (header.display || !displayOnly) {
			headerStrings = append(headerStrings, header.text)
		}
	}

	return headerStrings
}

func makeTable(headers []Header, contacts []Contact, click func(Contact)) *widget.Box {
	rows := ContactsToStrings(contacts)
	headerStrings := headersToText(headers, true)
	columnsAll := rowsToColumns(headers, rows)

	// remove the non-displaying columns
	var columns [][]string
	for idx, header := range headers {
		if (header.display) {
			columns = append(columns, columnsAll[idx])
		}
	}

	objects := make([]fyne.CanvasObject, len(columns)+1)
	for k, col := range columns {
		box := widget.NewVBox(widget.NewLabelWithStyle(headerStrings[k], fyne.TextAlignLeading, fyne.TextStyle{Bold: true}))
		for _, val := range col {
			box.Append(widget.NewLabel(val))
		}
		objects[k] = box
	}

	// add the edit buttons
	editbox := widget.NewVBox(widget.NewLabel(""))
	for i := range rows {
		row := i
		editbox.Append(widget.NewButton("Edit", func() {
			var contact Contact = contacts[row]
			fmt.Println("Id=" + strconv.Itoa(contact.ID))
			click(contact)
		}))
	}
	objects[len(columns)] = editbox

	return widget.NewHBox(objects...)
}

func rowsToColumns(headers []Header, rows [][]string) [][]string {
	columns := make([][]string, len(headers))
	for _, row := range rows {
		for colK := range row {
			columns[colK] = append(columns[colK], row[colK])
		}
	}
	return columns
}
