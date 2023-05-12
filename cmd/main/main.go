package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	bk "github.com/ilfey/Go-Phone-Book/internal/app/book"
	ai "github.com/ilfey/Go-Phone-Book/internal/pkg/ansi"
)

const (
	FILE_NAME = "dirictories.csv"
)

func FindContact(pb *bk.PhoneBook) (int, *bk.Contact) {
	var indeces []int

	ok := bk.GetConfirm("Вы хотите найти запись по имени контакта? (Y/n) >>> ", true)
	if ok {
		username := bk.GetUsername(pb.Scanner)
		indeces = pb.FindByUsername(username)
	} else {
		phone := bk.GetPhone(pb.Scanner)
		indeces = pb.FindByPhone(phone)
	}

	if len(indeces) == 0 {
		ai.Errorln("Контакт не найден.")
		return -1, nil
	}

	return bk.GetContact(indeces, pb)
}

func main() {
	pb, err := bk.NewPhoneBook(FILE_NAME)
	if err != nil {
		log.Fatal(err)
	}

	pb.AddCommands([]*bk.Command{
		{
			Title:       "print",
			Description: "Выводит список на печать.",
			Action: func(pb *bk.PhoneBook) error {
				bk.PrintlnContacts(pb.Contacts)
				return nil
			},
		},
		{
			Title:       "create",
			Description: "Создает новую запись.",
			Action: func(pb *bk.PhoneBook) error {
				contact := &bk.Contact{}

				contact.Username = bk.GetUsername(pb.Scanner)

				fmt.Print("Введите телефон >>> ")
				contact.Phone = bk.GetString() // TODO: add validation

				pb.Contacts = append(pb.Contacts, contact)

				ai.Successln("Запись создана.")

				return nil
			},
		},
		{
			Title:       "edit",
			Description: "Изменяет запись.",
			Action: func(pb *bk.PhoneBook) error {
				contactIndex, contact := FindContact(pb)
				if contactIndex < 0 {
					return nil
				}

				c := *contact

				for {
					bk.PrintlnContact(contact)

					fmt.Println("Введите новое имя. (Enter - не изменять)")

					username := bk.GetUsernameOrEmpty(pb.Scanner)
					if username != "" {
						c.Username = username
					}

					fmt.Println("Введите новый телефон (без знака \"+\"). (Enter - не изменять)")

					phone := bk.GetPhoneOrEmpty(pb.Scanner)
					if phone != "" {
						c.Phone = "+" + phone
					}

					if c.Username == contact.Username && c.Phone == contact.Phone {
						break
					}

					bk.PrintlnContact(&c)
					
					ok := bk.GetConfirm("Проверьте, все ли данные введены правильно. (Y/n) >>> ", true)
					if ok {
						pb.Contacts[contactIndex] = &c
						ai.Successln("Запись обновлена.")
						break
					}
				}
				return nil
			},
		},
		{
			Title:       "delete",
			Description: "Удаляет запись.",
			Action: func(pb *bk.PhoneBook) error {
				contactIndex, contact := FindContact(pb)
				if contactIndex < 0 {
					return nil
				}

				bk.PrintlnContact(contact)

				ok := bk.GetConfirm("Вы действительно хотите удалить эту запись? (Y/n) >>> ", true)
				if ok {
					pb.Contacts[contactIndex] = pb.Contacts[len(pb.Contacts)-1]
					pb.Contacts[len(pb.Contacts)-1] = nil
					pb.Contacts = pb.Contacts[:len(pb.Contacts)-1]
					ai.Successln("Запись удалена.")
				}

				return nil
			},
		},
		{
			Title:       "find",
			Description: "Поиск записи.",
			Action: func(pb *bk.PhoneBook) error {

				var indeces []int

				ok := bk.GetConfirm("Вы хотите найти запись по имени контакта? (Y/n) >>> ", true)
				if ok {
					username := bk.GetUsername(pb.Scanner)
					indeces = pb.FindByUsername(username)
				} else {
					phone := bk.GetPhone(pb.Scanner)
					indeces = pb.FindByPhone(phone)
				}

				if len(indeces) == 0 {
					ai.Errorln("Контакт не найден.")
					return nil
				}

				var contacts []*bk.Contact

				for _, i := range indeces {
					contacts = append(contacts, pb.Contacts[i])
				}

				bk.PrintlnContacts(contacts)

				return nil
			},
		},
	})

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGINT)

	go func() {
		<-signalChan
		err := pb.SavePhoneBook()
		if err != nil {
			ai.Errorln("Данные не удалось сохранить. Приносим свои глубочайшие извинения.")
		} else {
			ai.Successln("Данные сохранены.")
		}

		os.Exit(0)
	}()

	err = pb.Run()
	if err != nil {
		log.Fatal(err)
	}
}
