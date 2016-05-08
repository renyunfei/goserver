package main

import (
	. "conn"
	. "utils"
)

func main() {
	ParseConf(Conf)

	if err := InitServer(Conf); err != nil {
		return
	}

	Serv.Run()
}
