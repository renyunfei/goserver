package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	fname := "log/t.log"
	fin, err := os.Create(fname)
	defer fin.Close()
	if err != nil {
		fmt.Println("%s", err.Error())
		return
	}

	fname2 := "log/t1.log"
	fin2, err := os.Create(fname2)
	defer fin2.Close()
	if err != nil {
		fmt.Println("%s", err.Error())
		return
	}

	slog := log.New(fin, "[TEST]", log.Lshortfile)
	slog.Printf("%s", "Lshortfile")

	slog.SetFlags(log.Ldate)
	slog.Printf("%s", "log.Ldate")

	slog.SetFlags(log.Ltime)
	slog.Printf("%s", "log.Ltime")

	slog.SetOutput(fin2)

	slog.SetFlags(log.Lmicroseconds)
	slog.Printf("%s", "log.Lmicroseconds")

	slog.SetFlags(log.Llongfile)
	slog.Printf("%s", "log.Llongfile")

	slog.SetFlags(log.LUTC)
	slog.Printf("%s", "log.LUTC")

	slog.SetFlags(log.LstdFlags)
	slog.Printf("%s", "log.LstdFlags")

	s1 := "abcd"
	s2 := "aabcd"

	if s1 == s2 {
		fmt.Println("s1==s2")
	} else {
		fmt.Println("s1!=s2")
	}
}
