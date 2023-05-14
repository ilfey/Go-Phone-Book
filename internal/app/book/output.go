package book

import (
	"fmt"
	"strconv"

	ai "github.com/ilfey/Go-Phone-Book/internal/pkg/ansi"
)

// Вывести контакт
func PrintlnContact(c *Contact) {
	fmt.Println(ai.BOLD+ai.BG_WHITE+"               Имя                "+ai.MAGENTA+"#"+ai.NOCOLOR+"     Телефон    ", ai.NOSTYLE)
	fmt.Println(ai.BOLD + ai.MAGENTA + "==================================#=================" + ai.NOSTYLE)
	username := c.Username + "                                "[len(c.Username):]
	phone := c.Phone + "                "[len(c.Phone):]
	fmt.Println(ai.BOLD + " " + username + " " + ai.BOLD + ai.MAGENTA + "#" + ai.NOCOLOR + " " + phone + ai.NOSTYLE)
	fmt.Println(ai.BOLD + ai.MAGENTA + "==================================#=================" + ai.NOSTYLE)
}

// Вывести контакты
func PrintlnContacts(cs []*Contact) {
	fmt.Println(ai.BOLD+ai.BG_WHITE+"  №  "+ai.MAGENTA+"#"+ai.NOCOLOR+"               Имя                "+ai.MAGENTA+"#"+ai.NOCOLOR+"     Телефон    ", ai.NOSTYLE)
	fmt.Println(ai.BOLD + ai.MAGENTA + "=====#==================================#=================" + ai.NOSTYLE)

	for i, c := range cs {
		i++
		index := strconv.Itoa(i) + "   "[len(strconv.Itoa(i)):]
		username := c.Username + "                                "[len(c.Username):]
		phone := c.Phone + "                "[len(c.Phone):]
		if i%2 == 0 {
			fmt.Println(ai.BG_BLACK + " " + index + " " + ai.BOLD + ai.MAGENTA + "#" + ai.NOCOLOR + " " + username + " " + ai.BOLD + ai.MAGENTA + "#" + ai.NOCOLOR + " " + phone + ai.NOSTYLE)
			continue
		}
		fmt.Println(" " + index + " " + ai.BOLD + ai.MAGENTA + "#" + ai.NOCOLOR + " " + username + " " + ai.BOLD + ai.MAGENTA + "#" + ai.NOCOLOR + " " + phone + ai.NOSTYLE)
	}

	fmt.Println(ai.BOLD + ai.MAGENTA + "=====#==================================#=================" + ai.NOSTYLE)
}
