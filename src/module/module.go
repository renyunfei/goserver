package module

import (
	"errors"
	"fmt"
	"strings"
)

var ErrTest = errors.New("test init error")

func On_init() error {
	fmt.Println("on init")
	return nil
}

func On_recv(data []byte) ([]byte, int) {
	ups := strings.ToUpper(string(data))
	fmt.Printf("%s\n", ups)

	return []byte(ups), 0
}

func On_exit() error {
	fmt.Println("on exit")

	return nil
}
