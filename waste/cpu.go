package waste

import (
	"crypto/rand"
	"log"
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

		log.Printf("[CPU] Successfully wasted on %s", time.Now().Format(time.RFC3339))
		time.Sleep(interval)
	}
}

func runCPUWorker(workCh <-chan struct{}, doneCh chan<- struct{}) {
	buffer, cipher := newCPUBufferAndCipher()

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

func newCPUBufferAndCipher() ([]byte, *chacha20.Cipher) {
	var buffer []byte
	if len(Buffers) > 0 {
		buffer = make([]byte, 4*MiB)
		copy(buffer, Buffers[0].B[:4*MiB])
	} else {
		buffer = make([]byte, 4*MiB)
	}
	if _, err := rand.Read(buffer); err != nil {
		log.Panicf("failed to initialize CPU buffer: %v", err)
	}

	cipher, err := chacha20.NewUnauthenticatedCipher(buffer[:32], buffer[:24])
	if err != nil {
		log.Panicf("failed to initialize CPU cipher: %v", err)
	}
	return buffer, cipher
}
