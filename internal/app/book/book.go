package book

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"

	ai "github.com/ilfey/Go-Phone-Book/internal/pkg/ansi"
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
	errFileDamage = fmt.Errorf("файл с записями поврежден")
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
				Description: "Показывает справку.",
				Action: func(pb *PhoneBook) error {
					PrintHelp(pb.commands)
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

	var contacts []Contact

	defer func() {
		pb.Contacts = &contacts
	}()

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

func PrintlnContact(c *Contact) {
	fmt.Println(ai.BOLD + ai.INVERSE + "       Имя       " + ai.BG_MAGENTA + "#" + ai.NOSTYLE + ai.BOLD + ai.INVERSE + "     Телефон    " + ai.NOSTYLE)
	fmt.Println(ai.BOLD + ai.MAGENTA + "=================#================" + ai.NOSTYLE)
	username := c.Username + "                "[len(c.Username):]
	phone := c.Phone + "                "[len(c.Phone):]
	fmt.Println(username + " # " + phone)
	fmt.Println(ai.BOLD + ai.MAGENTA + "=================#================" + ai.NOSTYLE)
}

func PrintlnContacts(cs *[]Contact) {
	fmt.Println(ai.BOLD+ai.INVERSE+"  №  "+ai.BG_MAGENTA+"#"+ai.NOSTYLE+ai.BOLD+ai.INVERSE+"       Имя        "+ai.BG_MAGENTA+"#"+ai.NOSTYLE+ai.BOLD+ai.INVERSE+"    Телефон    ", ai.NOSTYLE)
	fmt.Println(ai.BOLD + ai.MAGENTA + "=====#==================#================" + ai.NOSTYLE)

	for i, c := range *cs {
		i++
		index := strconv.Itoa(i) + "   "[len(strconv.Itoa(i)):]
		username := c.Username + "                "[len(c.Username):]
		phone := c.Phone + "                "[len(c.Phone):]
		fmt.Println(" " + index + " " + ai.BOLD + ai.MAGENTA + "#" + ai.NOSTYLE + " " + username + " " + ai.BOLD + ai.MAGENTA + "#" + ai.NOSTYLE + " " + phone)
	}

	fmt.Println(ai.BOLD + ai.MAGENTA + "=====#==================#================" + ai.NOSTYLE)
}

func PrintHelp(cmds *[]Command) {
	for _, cmd := range *cmds {
		fmt.Println(ai.BG_GREEN + ai.BOLD + "Команда:" + ai.NOSTYLE + " " + ai.BOLD + ai.UNDERLINE + cmd.Title + ai.NOSTYLE)
		fmt.Println(ai.BG_MAGENTA + ai.BOLD + "Описание:" + ai.NOSTYLE + " " + ai.ITALIC + cmd.Description + ai.NOSTYLE)
		fmt.Println()
	}
}