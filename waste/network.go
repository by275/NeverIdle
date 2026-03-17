package waste

import (
	"context"
	"math/rand"
	"runtime"
	"time"

	"github.com/by275/neveridle/internal/log"
	"github.com/showwin/speedtest-go/speedtest"
)

const (
	networkFetchTimeout = 30 * time.Second
	networkPingTimeout  = 15 * time.Second
	networkTestTimeout  = 2 * time.Minute
)

func Network(interval time.Duration, connectionCount int) {
	cache := false
	speedtestClient := speedtest.New()
	speedtestClient.SetNThread(connectionCount)
	var targets speedtest.Servers
	for {
		if !cache {
			fetchCtx, cancelFetch := context.WithTimeout(context.Background(), networkFetchTimeout)
			_, err := speedtestClient.FetchUserInfoContext(fetchCtx)
			cancelFetch()
			if err != nil {
				log.Logf("NET", "Error when fetching user info: %v", err)
				sleepWithTimeout(time.Minute)
				continue
			}

			serverCtx, cancelServers := context.WithTimeout(context.Background(), networkFetchTimeout)
			serverList, err := speedtest.FetchServerListContext(serverCtx)
			cancelServers()
			if err != nil {
				log.Logf("NET", "Error when fetching servers: %v", err)
				sleepWithTimeout(time.Minute)
				continue
			}

			targets = *serverList.Available()
			if len(targets) == 0 {
				log.Logf("NET", "No available server to test, retrying in 5s")
				sleepWithTimeout(5 * time.Second)
				continue
			}
			if float64(len(targets))/float64(len(serverList)) > 0.5 {
				cache = true
			}
		}

		// pick random as main server
		s := targets[rand.Int31n(int32(len(targets)))]

		pingCtx, cancelPing := context.WithTimeout(context.Background(), networkPingTimeout)
		err := s.PingTestContext(pingCtx, nil)
		cancelPing()
		if err != nil {
			s.Latency = -1
		}

		downloadCtx, cancelDownload := context.WithTimeout(context.Background(), networkTestTimeout)
		err = s.MultiDownloadTestContext(downloadCtx, targets)
		cancelDownload()
		if err != nil {
			s.DLSpeed = -1
		}

		uploadCtx, cancelUpload := context.WithTimeout(context.Background(), networkTestTimeout)
		err = s.MultiUploadTestContext(uploadCtx, targets)
		cancelUpload()
		if err != nil {
			s.ULSpeed = -1
		}

		log.Logf("NET", "SpeedTest Ping=%s Download=%v Upload=%v mainServer=%s", s.Latency, s.DLSpeed, s.ULSpeed, s.String())

		speedtestClient.Manager.Reset()
		runtime.GC()
		sleepWithTimeout(interval)
	}
}

func sleepWithTimeout(d time.Duration) {
	timer := time.NewTimer(d)
	defer timer.Stop()
	<-timer.C
}
