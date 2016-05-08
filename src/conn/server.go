package conn

import (
	"container/list"
	"fmt"
	"log"
	"module"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	. "utils"
)

var (
	Serv = &Server{}
	Conf = &Configure{}
)

type Stat struct {
	gonum int
	alive int
}

type Server struct {
	listener *net.TCPListener
	clients  *list.List

	//Level int

	conf *Configure
	bind string
	stat Stat
	stop bool
	wg   sync.WaitGroup
	sig  chan os.Signal

	newConn chan Client
}

func InitServer(conf *Configure) error {
	SavePid()

	// Init Server
	Serv.conf = conf
	Serv.bind = conf.Bind
	Serv.newConn = make(chan Client, 100)
	Serv.clients = list.New()
	Serv.sig = make(chan os.Signal)

	// Set signal
	signal.Notify(Serv.sig, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT)
	signal.Ignore(syscall.SIGPIPE)

	// Listen bind
	addr, err := net.ResolveTCPAddr("tcp", Serv.bind)
	if err != nil {
		log.Fatal(err)
		return err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Println(err)
		return err
	}

	Serv.listener = l

	Serv.wg.Add(1)
	go func() {
		defer Serv.wg.Done()

		for {
			select {
			case c, ok := <-Serv.newConn:
				if ok {
					fmt.Println("new client")
					Serv.clients.PushBack(c)
				}
			case <-Serv.sig:
				Serv.stop = true

				//Notify listener not accept new connection
				t := time.Now().Add(time.Millisecond * 100)
				Serv.listener.SetDeadline(t)
				fmt.Fprintf(os.Stderr, "Receive stop signal, exiting...\n")

				//Notify client read close
				for c := Serv.clients.Front(); c != nil; c = c.Next() {
					c := c.Value.(Client)
					c.Req.Conn.SetDeadline(t)
				}

				return
			}
		}
	}()

	//process time event
	/*
		Serv.wg.Add(1)
		go func() {
			defer Serv.wg.Done()
			for {
				time.Sleep(time.Second)
			}
		}()
	*/
	return nil
}

func (self *Server) Run() {
	defer func() {
		err := module.On_exit()
		if err != nil {
			fmt.Println("On exit failed [%s], but we always exit because can not do nothing!", err.Error())
		}

		self.listener.Close()
		fmt.Println("Server exit")
	}()

	err := module.On_init()
	if err != nil {
		fmt.Println("%s", err.Error())
		return
	}

	for {
		if !self.stop {
			conn, err := self.listener.AcceptTCP()
			if err != nil {
				if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
					continue
				}

				fmt.Println(err)
			}

			self.wg.Add(1)
			go Handle_conn(conn)
		} else {
			self.wg.Wait()
			break
		}
	}
}
