package balePlugins

import (
	"net"

	"github.com/ALiwoto/ssg/ssg"
)

var (
	OwnersId []int64
	BotPairs []IdTuple
)

var BaleConnectionsPool = ssg.NewSafeMap[string, net.Conn]()
