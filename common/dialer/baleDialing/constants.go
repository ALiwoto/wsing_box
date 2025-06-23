package baleDialing

import "time"

const NetworkBale = "bale"

const (
	MaxGoRoutines = 50
	// MaxCharLen    = 4030
	MaxCharLen         = 4060
	MinCharLen         = 1020
	MaxDataSendRetries = 10
)

const (
	// delay per send per bot.
	DelayPerSend = time.Second

	bufferFlushInterval = 2 * time.Second
)

const (
	baleCommandCloseConn = "endMusical"
)

const (
	TooManyRequestStr = "Too Many Requests: retry after"
)
