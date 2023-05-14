package book

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"

	ai "github.com/ilfey/Go-Phone-Book/internal/pkg/ansi"
)

// Получить имя
func GetUsername(sc *bufio.Scanner, cmd *Command) string {
	reg := regexp.MustCompile(USERNAME_PATTERN)

	for {
		cmd.ScanCtxWithMsg("Введите имя")

		username, ok := GetLine(sc)
		if !ok {
			continue
		}

		if reg.MatchString(username) {
			return username
		}

		ai.Errorln("Данное имя не допустимо.")
	}
}

// Получить имя или пустую строку
func GetUsernameOrEmpty(sc *bufio.Scanner, cmd *Command) string {
	reg := regexp.MustCompile(USERNAME_PATTERN)

	for {
		cmd.ScanCtxWithMsg("Введите имя (Enter - пропустить)")

		username, ok := GetLine(sc)
		if !ok || username == "" {
			return ""
		}

		if reg.MatchString(username) {
			return username
		}

		ai.Errorln("Данное имя не допустимо.")
	}
}

// Получить телефон
func GetPhone(sc *bufio.Scanner, cmd *Command) string {
	reg := regexp.MustCompile(PHONE_PATTERN)

	for {
		cmd.ScanCtxWithMsg("Введите телефон (без \"+\")")

		phone, ok := GetLine(sc)
		if !ok {
			continue
		}

		phone = "+" + phone

		if reg.MatchString(phone) {
			return phone
		}

		ai.Errorln("Данный номер не допустим.")
	}
}

// Получить телефон или пустую строку
func GetPhoneOrEmpty(sc *bufio.Scanner, cmd *Command) string {
	reg := regexp.MustCompile(PHONE_PATTERN)

	for {
		cmd.ScanCtxWithMsg("Введите телефон (Enter - пропустить)")

		phone, ok := GetLine(sc)
		if !ok || phone == "" {
			return ""
		}

		phone = "+" + phone

		if reg.MatchString(phone) {
			return phone
		}

		ai.Errorln("Данный номер не допустим.")
	}
}

// Получить строку из консоли
func GetLine(sc *bufio.Scanner) (string, bool) {
	if sc.Scan() {
		input := sc.Text()
		return strings.TrimSpace(input), true
	}

	return "", false
}

// Получить строку до пробела из консоли
func GetString() string {
	var str string
	fmt.Scanln(&str)
	return str
}

// Получить число из консоли
func GetInt() int {
	var i int
	fmt.Scanln(&i)
	return i
}

// Получить от пользователя подтверждение
func GetConfirm(question string, cmd *Command, fallback bool) bool {
	cmd.ScanCtxWithMsg(question)

	var ok string
	fmt.Scanln(&ok)
	if fallback {
		if ok != "n" {
			return fallback
		}

		return false
	}

	if ok != "y" {
		return fallback
	}

	return true
}

// Требовать от пользователя выбрать контакт
func GetContact(indeces []int, pb *PhoneBook, cmd *Command) (int, *Contact) {
	var contactIndex int
	var contact *Contact

	if len(indeces) == 1 {
		contactIndex = indeces[0]
		contact = pb.Contacts[contactIndex]
	} else {
		var contacts []*Contact

		for _, i := range indeces {
			contacts = append(contacts, pb.Contacts[i])
		}

		PrintlnContacts(contacts)

		ai.Warnln("Было найдено несколько записей. Введите номер записи.")

		for {
			cmd.ScanCtx()

			index := GetInt()

			if index > len(indeces) || index < 1 {
				ai.Errorln("Пожалуйста, введите правильный номер записи.")
				continue
			}

			contactIndex = indeces[index-1]
			contact = pb.Contacts[contactIndex]
			break
		}
	}

	return contactIndex, contact
}
