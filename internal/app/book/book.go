package book

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	ai "github.com/ilfey/Go-Phone-Book/internal/pkg/ansi"
)

const (
	FILE_PERMISSION  = 0666
	USERNAME_PATTERN = `^[a-zA-Zа-яА-Я\s]{2,32}$`
	PHONE_PATTERN    = `[+]\d{8,15}$`
	FILE_PATTERN     = `^[a-zA-Zа-яА-Я\s]{2,32};[+]\d{8,15}$`
	SAVE_FORMAT      = "%s;%s\n"
)

var (
	errFileDamage = fmt.Errorf("файл с записями поврежден")
)

type Command struct {
	Title       string
	Description string
	Action      func(*PhoneBook, *Command) error
}

func (cmd *Command) PrintCtx() {
	fmt.Printf("%s[ %s ]%s ", ai.BG_BLUE, cmd.Title, ai.NORMAL)
}

func (cmd *Command) ScanCtx() {
	fmt.Printf("%s[ %s ]%s >>> ", ai.BG_BLUE, cmd.Title, ai.NORMAL)
}

func (cmd *Command) ScanCtxWithMsg(msg string) {
	fmt.Printf("%s[ %s ]%s %s >>> ", ai.BG_BLUE, cmd.Title, ai.NORMAL, msg)
}

type Contact struct {
	Username string
	Phone    string
}

type PhoneBook struct {
	Scanner  *bufio.Scanner
	Contacts []*Contact
	commands []*Command
	filename string
}

// Создать объект книги
func NewPhoneBook(f string) (*PhoneBook, error) {
	pb := &PhoneBook{
		Scanner: bufio.NewScanner(os.Stdin),
		commands: []*Command{
			{
				Title:       "help",
				Description: "Показывает справку.",
				Action: func(pb *PhoneBook, _ *Command) error {
					for _, cmd := range pb.commands {
						fmt.Println(ai.BG_GREEN + ai.BOLD + "Команда:" + ai.NOSTYLE + " " + ai.BOLD + ai.UNDERLINE + cmd.Title + ai.NOSTYLE)
						fmt.Println(ai.BG_MAGENTA + ai.BOLD + "Описание:" + ai.NOSTYLE + " " + ai.ITALIC + cmd.Description + ai.NOSTYLE)
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

// Добавить команды
func (pb *PhoneBook) AddCommands(cmds []*Command) {
	pb.commands = append(pb.commands, cmds...)
}

func (pb *PhoneBook) Run() error {
	for {
		fmt.Print(">>> ")

		cmdTitle := GetString()

		for _, cmd := range pb.commands {
			if cmdTitle == cmd.Title {
				cmd.Action(pb, cmd)
			}
		}
	}
}

// Загрузить данные
func (pb *PhoneBook) LoadPhoneBook() error {
	file, err := os.OpenFile(pb.filename, os.O_RDONLY|os.O_CREATE, FILE_PERMISSION)
	if err != nil {
		return err
	}

	defer file.Close()

	reg, err := regexp.Compile(FILE_PATTERN)
	if err != nil {
		return err
	}

	sc := bufio.NewScanner(file)

	var contacts []*Contact

	for sc.Scan() {
		text := sc.Text()

		if !reg.MatchString(text) {
			fmt.Println(text)
			return errFileDamage
		}

		vals := strings.Split(text, ";")

		contacts = append(contacts, &Contact{
			Username: vals[0],
			Phone:    vals[1],
		})
	}

	pb.Contacts = contacts
	return nil
}

// Сохранить книгу
func (pb *PhoneBook) SavePhoneBook() error {
	file, err := os.OpenFile(pb.filename, os.O_WRONLY|os.O_TRUNC, FILE_PERMISSION)
	if err != nil {
		return err
	}

	defer file.Close()

	wr := bufio.NewWriter(file)

	for _, c := range pb.Contacts {
		_, err = fmt.Fprintf(wr, SAVE_FORMAT, c.Username, c.Phone)
		if err != nil {
			return err
		}
	}

	return wr.Flush()
}

// Найти пользователя по имени
func (pb *PhoneBook) FindByUsername(username string) []int {
	var contacts []int

	for i, c := range pb.Contacts {
		if c.Username == username {
			contacts = append(contacts, i)
		}
	}

	return contacts
}

// Найти пользователя по телефону
func (pb *PhoneBook) FindByPhone(phone string) []int {
	var contacts []int

	if !strings.HasPrefix(phone, "+") {
		phone = "+" + phone
	}

	for i, c := range pb.Contacts {
		if c.Phone == phone {
			contacts = append(contacts, i)
		}
	}

	return contacts
}
