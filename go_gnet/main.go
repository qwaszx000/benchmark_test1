package main

import (
	"github.com/panjf2000/gnet/v2"
)

type TestServerGnet struct {
	gnet.BuiltinEventEngine //default implementations of gnet.EventHandler methods
}

func (server *TestServerGnet) OnBoot(eng gnet.Engine) gnet.Action {
	//log.Printf("Launched server on http://127.0.0.1:8080")
	return gnet.None
}

func (server *TestServerGnet) OnOpen(con gnet.Conn) ([]byte, gnet.Action) {
	//log.Printf("Connection from %s", con.RemoteAddr())
	return nil, gnet.None
}

func async_write_handler(con gnet.Conn, err error) error {
	if err != nil {
		//log.Printf("Error con.Next: %s", err)
		return err
	}

	con.Close()
	return nil
}

func (server *TestServerGnet) OnTraffic(con gnet.Conn) gnet.Action {
	resp_data := []byte("HTTP/1.1 200 OK\r\nServer: gnet\r\nContent-Type: text/plain\r\nContent-Length: 12\r\n\r\nHello world!")

	buff := make([]byte, 512)
	_, err := con.Read(buff)
	con.Discard(-1)
	if err != nil {
		//log.Printf("Error con.Next: %s\n", err)
		return gnet.Close
	}

	if is_data_ok(buff) {
		//log.Printf("Data ok\n")
		go con.AsyncWrite(resp_data, async_write_handler)
		return gnet.None
	}

	return gnet.Close
}

func main() {

	server := TestServerGnet{}
	gnet.Run(&server, "tcp://0.0.0.0:8080", gnet.WithMulticore(true))
}
