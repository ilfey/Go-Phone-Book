package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const (
	DB_FILE_NAME = "dirictories.json"
)

type Contact struct {
	Username string `json:"username"`
	Phone    string `json:"phone"`
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

// Загружает книгу из файла
// TODO: Create load to CSV
func (pb *PhoneBook) LoadPhoneBook() error {
	file, err := os.Open(pb.Filename)
	if err != nil {
		return err
	}

	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	var contacts []Contact

	if err := json.Unmarshal(bytes, &contacts); err != nil {
		return err
	}

	pb.Contacts = &contacts

	return nil
}

// Сохраняет книгу
// TODO: Create save to CSV
func (pb *PhoneBook) SavePhoneBook() error {
	j, err := json.Marshal(pb.Contacts)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(pb.Filename, j, 0644)
}

// Ищет пользователя
func (pb *PhoneBook) FindByUsername(u string) (*Contact, error) {
	var contact *Contact
	for _, c := range *pb.Contacts {
		if c.Username == u {
			contact = &c
		}
	}

	if len(contact.Username) == 0 {
		return nil, errors.New("contact not found")
	}

	return contact, nil
}

func main() {

	pb, err := NewPhoneBook(DB_FILE_NAME)
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

			break

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

			break

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

			break
		}
	}
}
