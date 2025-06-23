package dataPlugin

import (
	"sync"
	"time"

	"github.com/ALiwoto/ssg/ssg"
)

var (
	processedMessages = func() *ssg.SafeEMap[string, bool] {
		m := ssg.NewSafeEMap[string, bool]()
		m.SetExpiration(5 * time.Minute)
		m.SetInterval(10 * time.Minute)
		m.EnableChecking()

		return m
	}()
	processMessageLock = &sync.Mutex{}
)
