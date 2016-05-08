package utils

import (
	"errors"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Configure struct {
	Bind       string
	LogPath    string
	LogLevel   int
	LogFlag    int
	Debug      int
	LogMaxSize int64
	Rdeadline  int64
	Wdeadline  int64
	Rbuf       int
	Wbuf       int
	Maxgo      int
	Idleclose  int
}

func SavePid() {
	pid := os.Getpid()

	fin, err := os.Create("server.pid")
	defer fin.Close()
	if err != nil {
		log.Println(err)

		return
	}

	_, err = fin.WriteString(strconv.Itoa(pid))
	if err != nil {
		log.Println(err)
	}
}

func TypeConv(value string, ntype string) (reflect.Value, error) {
	if ntype == "string" {
		return reflect.ValueOf(value), nil
	} else if ntype == "time.Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "int" {
		i, err := strconv.Atoi(value)
		return reflect.ValueOf(i), err
	} else if ntype == "int8" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int8(i)), err
	} else if ntype == "int32" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int64(i)), err
	} else if ntype == "int64" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(i), err
	} else if ntype == "float32" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(float32(i)), err
	} else if ntype == "float64" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(i), err
	}

	return reflect.ValueOf(value), errors.New("Unknow type: " + ntype)
}

func ParseConf(conf *Configure) error {
	fin, err := os.Open("./server.conf")
	defer fin.Close()
	if err != nil {
		return err
	}

	//we suppose the length of configure file lt 4096
	buf := make([]byte, 4096)
	_, err = fin.Read(buf)
	if err != nil {
		return err
	}

	items := strings.Split(string(buf), "\n")

	for _, value := range items {
		if !strings.HasPrefix(value, "#") {

			kv := strings.Split(value, "=")
			if len(kv) != 2 {
				continue
			}

			k := strings.Title(strings.TrimSpace(kv[0]))
			v := strings.TrimSpace(kv[1])

			field := reflect.ValueOf(conf).Elem().FieldByName(k)

			if !field.IsValid() {
				return errors.New("Filed is NOT valid!")
			}

			if !field.CanSet() {
				return errors.New("Field can NOT set!")
			}

			val, err := TypeConv(v, field.Type().Name())
			if err != nil {
				return err
			}

			field.Set(val)
		}
	}

	return nil
}
