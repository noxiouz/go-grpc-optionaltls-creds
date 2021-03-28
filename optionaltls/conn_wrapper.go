package optionaltls

import (
	"bytes"
	"io"
	"net"
)

// WrappedConn Imitates MSG_PEEK behaviour
// Unlike net.Conn is not thread-safe
type WrappedConn struct {
	net.Conn

	// Reader for already peeked bytes
	peekedReader io.Reader
}

func NewWrappedConn(conn net.Conn, peeked []byte) net.Conn {
	var peekedReader = io.MultiReader(bytes.NewReader(peeked), conn)
	return &WrappedConn{
		Conn:         conn,
		peekedReader: peekedReader,
	}
}

func (wc *WrappedConn) Read(b []byte) (n int, err error) {
	return wc.peekedReader.Read(b)
}
