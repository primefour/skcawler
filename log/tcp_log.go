package log

import (
	"encoding/json"
	"io"
	"log"
	"net"
)

// implements LoggerInterface.
// it writes messages in keep-live tcp connection.
type TcpWriter struct {
	lg             *log.Logger
	innerWriter    io.WriteCloser
	ReconnectOnMsg bool   `json:"reconnectOnMsg"`
	Reconnect      bool   `json:"reconnect"`
	Net            string `json:"net"`
	Addr           string `json:"addr"`
	Level          int    `json:"level"`
}

// create new Logger returning as LoggerInterface.
func NewTcpLogger() LoggerInterface {
	logger := new(TcpWriter)
	logger.Level = LevelDebug
	return logger
}

// init connection writer with json config.
// json config only need key "level".
func (c *TcpWriter) Init(config map[string]interface{}) error {
	conf, err := json.Marshal(config)
	if err != nil {
		return err
	}
	return json.Unmarshal(conf, c)
}

// write message in connection.
// if connection is down, try to re-connect.
func (c *TcpWriter) WriteMsg(msg string, level int) error {
	if level > c.Level {
		return nil
	}
	if c.neddedConnectOnMsg() {
		err := c.connect()
		if err != nil {
			return err
		}
	}

	if c.ReconnectOnMsg {
		defer c.innerWriter.Close()
	}
	c.lg.Println(msg)
	return nil
}

// implementing method. empty.
func (c *ConnWriter) Flush() {

}

// destroy connection writer and close tcp listener.
func (c *ConnWriter) Destroy() {
	if c.innerWriter != nil {
		c.innerWriter.Close()
	}
}

func (c *TcpWriter) connect() error {
	if c.innerWriter != nil {
		c.innerWriter.Close()
		c.innerWriter = nil
	}

	conn, err := net.Dial(c.Net, c.Addr)
	if err != nil {
		return err
	}

	if tcpConn, ok := conn.(*net.TCPConn); ok {
		tcpConn.SetKeepAlive(true)
	}

	c.innerWriter = conn
	c.lg = log.New(conn, "", log.Ldate|log.Ltime)
	return nil
}

func (c *TcpWriter) neddedConnectOnMsg() bool {
	if c.Reconnect {
		c.Reconnect = false
		return true
	}

	if c.innerWriter == nil {
		return true
	}

	return c.ReconnectOnMsg
}

func init() {
	Register("TcpLogger", NewTcpLogger)
}
