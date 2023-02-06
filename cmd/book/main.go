package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	. "github.com/ilfey/Go-Phone-Book/internal/app/book"
)

var (
	errContactNotFound = fmt.Errorf("запись не найдена")
)

func getContact(indeces *[]int, pb *PhoneBook) (int, *Contact) {
	var contactIndex int
	var contact Contact

	if len(*indeces) == 1 {
		contactIndex = (*indeces)[0]
		contact = (*pb.Contacts)[contactIndex]
	} else {
		var contacts []Contact

		for _, i := range *indeces {
			contacts = append(contacts, (*pb.Contacts)[i])
		}

		PrintlnContacts(&contacts)

		fmt.Println("Было найдено несколько записей с таким именем. Введите номер записи.")

		for {
			fmt.Print(">>> ")

			var index int
			fmt.Scanln(&index)

			if index > len(*indeces) || index < 1 {
				fmt.Println("Пожалуйста, введите правильный номер записи.")
				continue
			}

			contactIndex = (*indeces)[index-1]
			contact = (*pb.Contacts)[contactIndex]
			break
		}
	}

	return contactIndex, &contact
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
				PrintlnContacts(pb.Contacts)
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

				contactIndex, contact := getContact(contactIndeces, pb)

				c := *contact

				for {
					PrintlnContact(contact)
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
		{
			Title:       "delete",
			Description: "Удаляет запись.",
			Action: func(pb *PhoneBook) error {
				fmt.Print("Введите имя контакта >>> ")
				var u string
				fmt.Scanln(&u) // TODO: add validation

				contactIndeces := pb.FindByUsername(u)
				if contactIndeces == nil {
					fmt.Println(errContactNotFound)
					return nil
				}

				contactIndex, contact := getContact(contactIndeces, pb)

				PrintlnContact(contact)

				fmt.Print("Вы действительно хотите удалить эту запись? (Y/n) >>> ")
				var ok string
				fmt.Scanln(&ok)
				if ok != "n" {
					(*pb.Contacts)[contactIndex] = (*pb.Contacts)[len(*pb.Contacts)-1]
					// (*pb.Contacts)[len(*pb.Contacts)-1] = nil // TODO
					(*pb.Contacts) = (*pb.Contacts)[:len(*pb.Contacts)-1]
					fmt.Println("Запись удалена.")
				}

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
