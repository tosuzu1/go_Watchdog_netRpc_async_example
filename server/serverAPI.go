package main

import (
	"sync"
	"time"
)

type ServerAPI int

var watchDogElaspedTime time.Duration
var isWatchDogTiemout bool
var watchDogTimeout time.Duration
var timePast time.Time
var timeNow time.Time
var watchDogTimeLock sync.Mutex

func (s *ServerAPI) Init(requestTime_in_millisecond int) {
	watchDogTimeout = time.Duration(requestTime_in_millisecond) * time.Millisecond
	watchDogElaspedTime = watchDogTimeout
	isWatchDogTiemout = false
}

func (s *ServerAPI) StartWatchDog() {
	timePast = time.Now()

	for !isWatchDogTiemout {
		// decrement watchDog until
		timeNow = time.Now()
		timeDelta := timePast.Sub(timeNow)
		//log.Printf("timeDelta %v\n", timeDelta)

		watchDogTimeLock.Lock()
		watchDogElaspedTime += timeDelta

		if 0 > watchDogElaspedTime {
			isWatchDogTiemout = true
		}
		watchDogTimeLock.Unlock()
		timePast = timeNow
	}
}

func (s *ServerAPI) IsWatchTimeout() bool {
	return isWatchDogTiemout
}

func (s *ServerAPI) getElapsedTime() time.Duration {
	return watchDogElaspedTime
}

func (s *ServerAPI) ClientHeartBeat(_ int, resp *int) error {
	watchDogTimeLock.Lock()
	watchDogElaspedTime = watchDogTimeout
	watchDogTimeLock.Unlock()

	return nil
}
