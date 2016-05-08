package conn

import (
	"encoding/gob"
	"errors"
	"io"
	"module"
	"net"
	"time"
)

var (
	ErrPackageLength = errors.New("Decode: Length is NOT equal the length of receive package")
)

const (
	least = 51
)

const (
	ESTABLISHED = iota
	CLOSED
	RETRY
)

type Request struct {
	addr net.Addr
	Conn *net.TCPConn

	state       int32
	timeout     time.Duration
	retry_count int32
}

type Handler interface {
	on_open() error
	on_recv([]byte) ([]byte, error)
	on_send() error
	on_close() error
	on_time() error
}

type Client struct {
	smsg Msg
	rmsg Msg

	msg chan int

	dec *gob.Decoder
	enc *gob.Encoder

	Req  *Request
	serv *Server
}

type Msg struct {
	Length uint32
	Cmd    uint32
	Ver    uint8
	SyncNo uint64
	Auth   [32]byte
	ErrNo  uint16
	Data   []byte
}

type Ticker struct {
	t uint64
}

func (c *Client) pack() error {
	c.smsg.Length = uint32(len(c.smsg.Data)) + 51
	c.smsg.SyncNo += 1

	return c.enc.Encode(c.smsg)
}

func (c *Client) unpack() error {
	c.rmsg = Msg{}
	err := c.dec.Decode(&c.rmsg)
	if err != nil && err != io.EOF {
		return err
	}

	if c.rmsg.Length != (51 + uint32(len(c.rmsg.Data))) {
		return ErrPackageLength
	}

	return nil
}

func NewClient(conn *net.TCPConn) *Client {
	client := new(Client)
	client.serv = Serv
	client.dec = gob.NewDecoder(conn)
	client.enc = gob.NewEncoder(conn)

	req := new(Request)
	req.addr = conn.RemoteAddr()
	req.state = ESTABLISHED
	req.Conn = conn

	client.Req = req

	return client
}

func Handle_conn(conn *net.TCPConn) {
	defer func() {
		//Log.Info("client from %s break", c.Req.addr.String())
		Serv.wg.Done()
		conn.Close()
	}()

	c := NewClient(conn)
	Serv.newConn <- *c

	Log.Info("new client from %s", c.Req.addr.String())

	for {
		err := c.unpack()
		if err != nil {
			Log.Warn("%s", err.Error())
			return
		}

		Log.Info("Cmd:%d Length:%d Data:%s Stop:%t Size:%d\n",
			c.rmsg.Cmd, c.rmsg.Length, c.rmsg.Data, Serv.stop, len(c.rmsg.Data))

		result, errno := module.On_recv(c.rmsg.Data)
		if errno != 0 {
			c.smsg.ErrNo = uint16(errno)
		} else {
			c.smsg.Data = result
		}

		c.pack()

		if Serv.stop {
			break
		}
	}
}
