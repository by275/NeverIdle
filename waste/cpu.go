package waste

import (
	"crypto/rand"
	"fmt"
	"time"

	"golang.org/x/crypto/chacha20"
)

const cpuWorkerCount = 8

func CPU(interval time.Duration) {
	workCh := make(chan struct{})
	doneCh := make(chan struct{}, cpuWorkerCount)
	for i := 0; i < cpuWorkerCount; i++ {
		go runCPUWorker(workCh, doneCh)
	}

	for {
		for i := 0; i < cpuWorkerCount; i++ {
			workCh <- struct{}{}
		}
		for i := 0; i < cpuWorkerCount; i++ {
			<-doneCh
		}

		fmt.Println("[CPU] Successfully wasted on", time.Now())
		time.Sleep(interval)
	}
}

func runCPUWorker(workCh <-chan struct{}, doneCh chan<- struct{}) {
	var buffer []byte
	if len(Buffers) > 0 {
		buffer = make([]byte, 4*MiB)
		copy(buffer, Buffers[0].B[:4*MiB])
	} else {
		buffer = make([]byte, 4*MiB)
	}
	_, _ = rand.Read(buffer)

	cipher, err := chacha20.NewUnauthenticatedCipher(buffer[:32], buffer[:24])
	if err != nil {
		panic(err)
	}

	for range workCh {
		for i := 0; i < 64; i++ {
			cipher.XORKeyStream(buffer, buffer)
		}

		newCipher, err := chacha20.NewUnauthenticatedCipher(buffer[:32], buffer[:24])
		if err == nil {
			cipher = newCipher
		}

		doneCh <- struct{}{}
	}
}
