// Package tcpkeepalive implements additional TCP keepalive control beyond what
// is currently offered by the net pkg.
//
// Only Linux >= 2.4 and OS X >= 10.8 are supported at this point, but
// patches for additional platforms are welcome.
package tcpkeepalive

import (
	"fmt"
	"net"

	"time"
)

// EnableKeepAlive enables TCP keepalive for the given conn, which must be a
// *tcp.TCPConn. The returned Conn allows overwriting the default keepalive
// parameters used by the operating system.
func EnableKeepAlive(conn net.Conn) (*Conn, error) {
	tcp, ok := conn.(*net.TCPConn)
	if !ok {
		return nil, fmt.Errorf("Bad conn type: %T", conn)
	}
	if err := tcp.SetKeepAlive(true); err != nil {
		return nil, err
	}
	file, err := tcp.File()
	if err != nil {
		return nil, err
	}
	fd := int(file.Fd())
	return &Conn{TCPConn: tcp, fd: fd}, nil
}

// Conn adds additional TCP keepalive control to a *net.TCPConn.
type Conn struct {
	*net.TCPConn
	fd int
}

// The time (in seconds) the connection needs to remain idle before TCP starts
// sending keepalive probes.
func (c *Conn) SetIdle(d time.Duration) error {
	return setIdle(c.fd, secs(d))
}

// The maximum number of keepalive probes TCP should send before dropping the
// connection.
func (c *Conn) SetCount(n int) error {
	return setCount(c.fd, n)
}

// The time (in seconds) between individual keepalive probes.
func (c *Conn) SetInterval(d time.Duration) error {
	return setInterval(c.fd, secs(d))
}

func secs(d time.Duration) int {
	d += (time.Second - time.Nanosecond)
	return int(d.Seconds())
}
