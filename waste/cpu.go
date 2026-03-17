package waste

import (
	"crypto/rand"
	"time"

	"github.com/by275/neveridle/internal/log"
	"golang.org/x/crypto/chacha20"
)

const cpuWorkerCount = 8

func CPU(interval time.Duration) {
	workCh := make(chan struct{})
	doneCh := make(chan struct{}, cpuWorkerCount)
	for range cpuWorkerCount {
		go runCPUWorker(workCh, doneCh)
	}

	for {
		for range cpuWorkerCount {
			workCh <- struct{}{}
		}
		for range cpuWorkerCount {
			<-doneCh
		}

		log.Logf("CPU", "Successfully wasted on %s", time.Now().Format(time.RFC3339))
		time.Sleep(interval)
	}
}

func runCPUWorker(workCh <-chan struct{}, doneCh chan<- struct{}) {
	buffer, cipher := newCPUBufferAndCipher()

	for range workCh {
		for range 64 {
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
	buffer := allocatedMemory.firstBufferPrefix(4 * MiB)
	if len(buffer) == 0 {
		buffer = make([]byte, 4*MiB)
	}
	if _, err := rand.Read(buffer); err != nil {
		log.Panicf("CPU", "failed to initialize CPU buffer: %v", err)
	}

	cipher, err := chacha20.NewUnauthenticatedCipher(buffer[:32], buffer[:24])
	if err != nil {
		log.Panicf("CPU", "failed to initialize CPU cipher: %v", err)
	}
	return buffer, cipher
}
