// Package tcpkeepalive implements additional TCP keepalive control beyond what
// is currently offered by the net pkg.
//
// Only Linux >= 2.4, DragonFly, FreeBSD, NetBSD and OS X >= 10.8 are supported
// at this point, but patches for additional platforms are welcome.
//
// See also: http://felixge.de/2014/08/26/tcp-keepalive-with-golang.html
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

// SetKeepAliveIdle sets the time (in seconds) the connection needs to remain
// idle before TCP starts sending keepalive probes.
func (c *Conn) SetKeepAliveIdle(d time.Duration) error {
	return setIdle(c.fd, secs(d))
}

// SetKeepAliveCount sets the maximum number of keepalive probes TCP should
// send before dropping the connection.
func (c *Conn) SetKeepAliveCount(n int) error {
	return setCount(c.fd, n)
}

// SetKeepAliveInterval sets the time (in seconds) between individual keepalive
// probes.
func (c *Conn) SetKeepAliveInterval(d time.Duration) error {
	return setInterval(c.fd, secs(d))
}

func secs(d time.Duration) int {
	d += (time.Second - time.Nanosecond)
	return int(d.Seconds())
}
