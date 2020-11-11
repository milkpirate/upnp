// +build windows

package upnp

import (
	"net"
)

func setTTL(conn *net.UDPConn, ttl int) error {
	return nil
}
