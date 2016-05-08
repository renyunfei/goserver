package main

import (
	. "conn"
	"fmt"
	"os"
	. "utils"
)

func main() {
	ParseConf(Conf)

	if err := InitServer(Conf); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return
	}

	Serv.Run()
}
