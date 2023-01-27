package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
)

const (
	FILE_NAME        = "dirictories.csv"
	FILE_PERMISSION  = 0666
	FILE_PATTERN     = `([A-z0-9])+[^;\s]`
	USERNAME_PATTERN = `` // TODO: add validation pattern
	PHONE_PATTERN    = `` // TODO: add validation pattern
	SAVE_FORMAT      = "%s;%s\n"
)

var (
	errFileDamage      = fmt.Errorf("the file is damaged")
	errContactNotFound = fmt.Errorf("contact not found")
)

type Command struct {
	Title       string
	Description string
	Action      func(*PhoneBook, *bufio.Scanner) error
}

type Contact struct {
	Username string
	Phone    string
}

type PhoneBook struct {
	Contacts *[]Contact
	commands *[]Command
	filename string
}

func NewPhoneBook(f string) (*PhoneBook, error) {
	pb := &PhoneBook{
		commands: &[]Command{
			{
				Title:       "help",
				Description: "Показывает справку",
				Action: func(pb *PhoneBook, s *bufio.Scanner) error {
					for _, cmd := range *pb.commands {
						fmt.Println("Команда:", cmd.Title)
						fmt.Println("Описание:", cmd.Description)
						fmt.Println()
					}

					return nil
				},
			},
		},
		filename: f,
	}

	err := pb.LoadPhoneBook()
	if err != nil {
		return nil, err
	}

	return pb, nil
}

func (pb *PhoneBook) AddCommands(cmds *[]Command) {
	newCmds := append(*((*pb).commands), *cmds...)
	pb.commands = &newCmds
}

func (pb *PhoneBook) Run() error {
	sc := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print(">>> ")
		sc.Scan()
		cmdTitle := sc.Text()

		for _, cmd := range *pb.commands {
			if cmdTitle == cmd.Title {
				cmd.Action(pb, sc)
			}
		}
	}
}

// Загрузить данные
func (pb *PhoneBook) LoadPhoneBook() error {
	file, err := os.OpenFile(pb.filename, os.O_RDWR|os.O_CREATE, FILE_PERMISSION)
	if err != nil {
		return err
	}

	defer file.Close()

	reg, err := regexp.Compile(FILE_PATTERN)
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
	file, err := os.OpenFile(pb.filename, os.O_RDWR, FILE_PERMISSION)
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
func (pb *PhoneBook) FindByUsername(contact *Contact) error {
	for _, c := range *pb.Contacts {
		if c.Username == contact.Username {
			contact = &c
			break
		}
	}

	if len(contact.Phone) == 0 {
		return errContactNotFound
	}

	return nil
}

func main() {
	pb, err := NewPhoneBook(FILE_NAME)
	if err != nil {
		log.Fatal(err)
	}

	pb.AddCommands(&[]Command{
		{
			Title:       "print",
			Description: "Выводит список на печать.",
			Action: func(pb *PhoneBook, s *bufio.Scanner) error {
				fmt.Println("Имя контакта    |Телефон")

				for _, el := range *pb.Contacts {
					fmt.Println(el.Username, el.Phone)
				}

				return nil
			},
		},
		{
			Title:       "create",
			Description: "Создает новую запись.",
			Action: func(pb *PhoneBook, s *bufio.Scanner) error {
				var contact Contact
				fmt.Print("Введите имя контакта >>> ")
				s.Scan()
				contact.Username = s.Text() // TODO: add validation

				fmt.Print("Введите телефон >>> ")
				s.Scan()
				contact.Phone = s.Text() // TODO: add validation

				*pb.Contacts = append(*pb.Contacts, contact)

				pb.SavePhoneBook() // FIXME: remove it

				return nil
			},
		},
		{
			Title:       "edit",
			Description: "Изменяет запись.",
			Action: func(pb *PhoneBook, s *bufio.Scanner) error {
				var contact Contact
				fmt.Print("Введите имя контакта >>> ")
				s.Scan()
				contact.Username = s.Text() // TODO: add validation

				err := pb.FindByUsername(&contact)
				if err != nil {
					fmt.Println(err)
					return nil
				}

				fmt.Println("Username:", contact.Username, "phone:", contact.Phone)

				fmt.Print("Введите новое имя (Enter - не изменять) >>> ")
				s.Scan()
				username := s.Text()
				fmt.Println(username)
				if username != "" {
					contact.Username = username
				}

				fmt.Print("Введите новый телефон (Enter - не изменять >>> ")
				s.Scan()
				phone := s.Text()
				if phone != "" {
					contact.Phone = phone
				}

				return nil
			},
		},
	})

	err = pb.Run()
	if err != nil {
		log.Fatal(err)
	}
}
