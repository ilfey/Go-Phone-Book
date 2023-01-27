package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
)

const (
	FILE_NAME          = "dirictories.csv"
	FILE_PERMISSION    = 0666
	REGULAR_EXPRESSION = `([A-z0-9])+[^;\s]`
	SAVE_FORMAT        = "%s;%s\n"
)

var (
	errFileDamage      = fmt.Errorf("the file is damaged")
	errContactNotFound = fmt.Errorf("contact not found")
)

type Contact struct {
	Username string
	Phone    string
}

type PhoneBook struct {
	Contacts *[]Contact
	Filename string
}

func NewPhoneBook(f string) (*PhoneBook, error) {
	pb := &PhoneBook{
		Filename: f,
	}

	err := pb.LoadPhoneBook()
	if err != nil {
		return nil, err
	}

	return pb, nil
}

// Загрузить данные
func (pb *PhoneBook) LoadPhoneBook() error {
	file, err := os.OpenFile(pb.Filename, os.O_RDWR|os.O_CREATE, FILE_PERMISSION)
	if err != nil {
		return err
	}

	defer file.Close()

	reg, err := regexp.Compile(REGULAR_EXPRESSION)
	if err != nil {
		return err
	}

	sc := bufio.NewScanner(file)

	if sc.Scan() {
		if !reg.MatchString(sc.Text()) {
			return errFileDamage
		}
	} else {
		return nil
	}

	var contacts []Contact

	for sc.Scan() {
		if !reg.MatchString(sc.Text()) {
			return errFileDamage
		}

		vals := reg.FindAllString(sc.Text(), -1)

		contacts = append(contacts, Contact{
			Username: vals[0],
			Phone:    vals[1],
		})
	}

	pb.Contacts = &contacts

	return nil
}

// Сохранить книгу
func (pb *PhoneBook) SavePhoneBook() error {
	file, err := os.OpenFile(pb.Filename, os.O_RDWR, FILE_PERMISSION)
	if err != nil {
		return err
	}

	defer file.Close()

	wr := bufio.NewWriter(file)

	if _, err := fmt.Fprintf(wr, SAVE_FORMAT, "username", "phone"); err != nil {
		return err
	}

	for _, c := range *pb.Contacts {
		_, err = fmt.Fprintf(wr, SAVE_FORMAT, c.Username, c.Phone)
		if err != nil {
			return err
		}
	}

	return wr.Flush()
}

// Найти пользователя
func (pb *PhoneBook) FindByUsername(u string) (*Contact, error) {
	var contact *Contact
	for _, c := range *pb.Contacts {
		if c.Username == u {
			contact = &c
		}
	}

	if len(contact.Username) == 0 {
		return nil, errContactNotFound
	}

	return contact, nil
}

func main() {

	pb, err := NewPhoneBook(FILE_NAME)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print(">>> ")
		scanner.Scan()
		action := scanner.Text()

		switch action {
		case "print":
			for _, el := range *pb.Contacts {
				fmt.Println("Username:", el.Username, "phone:", el.Phone)
			}

		case "create":
			var contact Contact
			fmt.Print("Enter username >>> ")
			scanner.Scan()
			contact.Username = scanner.Text()

			fmt.Print("Enter phone >>> ")
			scanner.Scan()
			contact.Phone = scanner.Text()

			*pb.Contacts = append(*pb.Contacts, contact)

			pb.SavePhoneBook()

		case "edit":
			fmt.Print("Enter username >>> ")
			scanner.Scan()
			username := scanner.Text()

			contact, err := pb.FindByUsername(username)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println("Username:", contact.Username, "phone:", contact.Phone)

				fmt.Print("Enter new username (Enter if not change) >>> ")
				scanner.Scan()
				username := scanner.Text()
				fmt.Println(username)
				if username != "" {
					contact.Username = username
				}

				fmt.Print("Enter new phone (Enter if not change) >>> ")
				scanner.Scan()
				phone := scanner.Text()
				if phone != "" {
					contact.Phone = phone
				}

			}

		case "save":
			err := pb.SavePhoneBook()
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
