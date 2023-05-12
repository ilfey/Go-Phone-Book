package book

import (
	"fmt"
	"strconv"

	ai "github.com/ilfey/Go-Phone-Book/internal/pkg/ansi"
)

func PrintlnContact(c *Contact) {
	fmt.Println(ai.BOLD + ai.INVERSE + "       Имя       " + ai.BG_MAGENTA + "#" + ai.NOSTYLE + ai.BOLD + ai.INVERSE + "     Телефон    " + ai.NOSTYLE)
	fmt.Println(ai.BOLD + ai.MAGENTA + "=================#================" + ai.NOSTYLE)
	username := c.Username + "                "[len(c.Username):]
	phone := c.Phone + "                "[len(c.Phone):]
	fmt.Println(username + " " + ai.BG_MAGENTA + "#" + ai.NOSTYLE + " " + phone)
	fmt.Println(ai.BOLD + ai.MAGENTA + "=================#================" + ai.NOSTYLE)
}

func PrintlnContacts(cs []*Contact) {
	fmt.Println(ai.BOLD+ai.INVERSE+"  №  "+ai.BG_MAGENTA+"#"+ai.NOSTYLE+ai.BOLD+ai.INVERSE+"       Имя        "+ai.BG_MAGENTA+"#"+ai.NOSTYLE+ai.BOLD+ai.INVERSE+"    Телефон    ", ai.NOSTYLE)
	fmt.Println(ai.BOLD + ai.MAGENTA + "=====#==================#================" + ai.NOSTYLE)

	for i, c := range cs {
		i++
		index := strconv.Itoa(i) + "   "[len(strconv.Itoa(i)):]
		username := c.Username + "                "[len(c.Username):]
		phone := c.Phone + "                "[len(c.Phone):]
		if i % 2 == 0 {
			fmt.Println(ai.BG_BLACK + " " + index + " " + ai.BOLD + ai.MAGENTA + "#" + ai.NOCOLOR + " " + username + " " + ai.BOLD + ai.MAGENTA + "#" + ai.NOCOLOR + " " + phone + ai.NOSTYLE)
			continue
		}
		fmt.Println(" " + index + " " + ai.BOLD + ai.MAGENTA + "#" + ai.NOCOLOR + " " + username + " " + ai.BOLD + ai.MAGENTA + "#" + ai.NOCOLOR + " " + phone + ai.NOSTYLE)
	}

	fmt.Println(ai.BOLD + ai.MAGENTA + "=====#==================#================" + ai.NOSTYLE)
}
