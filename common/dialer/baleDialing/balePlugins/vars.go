package balePlugins

import (
	"net"

	"github.com/ALiwoto/ssg/ssg"
)

var (
	OwnersId []int64
	BotPairs []PairsMinimalInfo
)

var BaleConnectionsPool = ssg.NewSafeMap[string, net.Conn]()

var HandleNewBaleConn func(connId string) error
