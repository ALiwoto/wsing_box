package baleDialing

import (
	"net"
	"sync"
	"sync/atomic"

	"github.com/ALiwoto/ssg/ssg"
	"github.com/sagernet/sing-box/common/dialer/baleDialing/gotgbot"
	"github.com/sagernet/sing-box/common/dialer/baleDialing/gotgbot/ext"
	"github.com/sagernet/sing-box/option"
)

//---------------------------------------------------------

type BaleBotContainer struct {
	Updater    *ext.Updater
	Dispatcher *ext.Dispatcher
	Bot        *gotgbot.Bot
	MyConfig   *option.BaleBotInfo
	BaleConfig *option.BaleConfiguration
	lock       *sync.Mutex
}

type BaleBotPair struct {
	InsideBot  *BaleBotContainer
	OutsideBot *BaleBotContainer
}

type BaleBotPairsContainer struct {
	// can add some other options here later maybe
	Pairs []*BaleBotPair
}

type BaleDialerContainer struct {
	Bots     *BaleBotPairsContainer
	connPool *ssg.SafeMap[string, BaleConn]
}

//---------------------------------------------------------

type BaleConn struct {
	ConnectionId string
	Address      *BaleFakeAddr
	pipe         net.Conn
	botsPool     []*BaleBotContainer
	writeLock    *sync.Mutex
	readLock     *sync.Mutex
	bufferedData []byte
	isClosed     *atomic.Bool
}

type PacketConn struct {
}

type BaleFakeAddr struct {
	NetworkStr string
	AddressStr string
}

//---------------------------------------------------------
