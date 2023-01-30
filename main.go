package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"syscall"
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
	Action      func(*PhoneBook) error
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
				Action: func(pb *PhoneBook) error {
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
	for {
		fmt.Print(">>> ")
		var cmdTitle string
		fmt.Scanln(&cmdTitle)

		for _, cmd := range *pb.commands {
			if cmdTitle == cmd.Title {
				cmd.Action(pb)
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
func (pb *PhoneBook) FindByUsername(username string) *[]int {
	var contacts []int

	for i, c := range *pb.Contacts {
		if c.Username == username {
			contacts = append(contacts, i)
		}
	}

	if len(contacts) == 0 {
		return nil
	}

	return &contacts
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
			Action: func(pb *PhoneBook) error {
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
			Action: func(pb *PhoneBook) error {
				var contact Contact
				fmt.Print("Введите имя контакта >>> ")
				fmt.Scanln(&contact.Username) // TODO: add validation

				fmt.Print("Введите телефон >>> ")
				fmt.Scanln(&contact.Phone) // TODO: add validation

				*pb.Contacts = append(*pb.Contacts, contact)

				fmt.Println("Запись создана.")

				return nil
			},
		},
		{
			Title:       "edit",
			Description: "Изменяет запись.",
			Action: func(pb *PhoneBook) error {
				fmt.Print("Введите имя контакта >>> ")
				var username string
				fmt.Scanln(&username) // TODO: add validation

				contactIndeces := pb.FindByUsername(username)
				if contactIndeces == nil {
					fmt.Println(errContactNotFound)
					return nil
				}

				var contactIndex int
				var contact Contact

				if len(*contactIndeces) == 1 {
					contactIndex = (*contactIndeces)[0]
					contact = (*pb.Contacts)[contactIndex]
				} else {
					fmt.Println("Было найдено несколько записей с таким именем. Введите номер записи, которую нужно изменить.")
					for i, el := range *contactIndeces {
						fmt.Printf("%d. Имя: %s Телефон: %s\n", i+1, (*pb.Contacts)[el].Username, (*pb.Contacts)[el].Phone)
					}

					for {
						var index int
						fmt.Scanln(&index)

						if index > len(*contactIndeces) || index < 1 {
							fmt.Println("Пожалуйста, введите правильный номер записи.")
							continue
						}

						contactIndex = (*contactIndeces)[index-1]
						contact = (*pb.Contacts)[contactIndex]
						break
					}
				}

				c := contact

				for {
					fmt.Println("Имя:", contact.Username, "Телефон:", contact.Phone)
					fmt.Print("Введите новое имя. (Enter - не изменять) >>> ")

					var username string
					fmt.Scanln(&username) // TODO: add validation
					if username != "" {
						c.Username = username
					}

					fmt.Print("Введите новый телефон. (Enter - не изменять) >>> ")
					var phone string
					fmt.Scanln(&phone) // TODO: add validation
					if phone != "" {
						c.Phone = phone
					}

					if c.Username == contact.Username && c.Phone == contact.Phone {
						break
					}

					fmt.Print("Проверьте, все ли данные введены правильно. (Y/n) >>> ")
					var ok string
					fmt.Scanln(&ok)
					if ok != "n" {
						(*pb.Contacts)[contactIndex] = c
						fmt.Println("Запись обновлена.")
						break
					}
				}
				return nil
			},
		},
	})

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-signalChan
		err := pb.SavePhoneBook()
		if err != nil {
			fmt.Println("Данные не сохранены.")
		} else {
			fmt.Println("Данные сохранены.")
		}

		os.Exit(0)
	}()

	err = pb.Run()
	if err != nil {
		log.Fatal(err)
	}
}
