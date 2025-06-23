package baleDialing

import (
	"net"

	"github.com/ALiwoto/ssg/ssg"
)

var cmdPrefixes = []rune{
	'/', '!',
}

var (
	baleConnectionsPool = ssg.NewSafeMap[string, net.Conn]()
)
