package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"

	"github.com/by275/neveridle/controller"
	"github.com/by275/neveridle/waste"
	"github.com/shirou/gopsutil/v3/mem"
)

const Version = "0.2.3"

const minMemoryReserveGiB = 1

var (
	FlagCPUPercent             = flag.Float64("cp", 0, "Target CPU waste ratio between 0 and 1")
	FlagCPU                    = flag.Duration("c", 0, "Interval for CPU waste")
	FlagMemory                 = flag.Int("m", 0, "GiB of memory waste")
	FlagNetwork                = flag.Duration("n", 0, "Interval for network speed test")
	FlagNetworkConnectionCount = flag.Int("t", 10, "Set concurrent connections for network speed test")
	FlagPriority               = flag.Int("p", 666, "Set process priority value")
)

func main() {
	printBanner()

	if err := parseAndValidateFlags(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		flag.PrintDefaults()
		os.Exit(2)
	}

	nothingEnabled := true
	applyPriority()
	nothingEnabled = startMemoryWaste(nothingEnabled)
	nothingEnabled = startCPUWaste(nothingEnabled)
	nothingEnabled = startNetworkWaste(nothingEnabled)

	if nothingEnabled {
		flag.PrintDefaults()
		return
	}

	waitForShutdown()
}

func printBanner() {
	fmt.Println("NeverIdle", Version, "- Getting worse from here.")
	fmt.Println("Platform:", runtime.GOOS, ",", runtime.GOARCH, ",", runtime.Version())
	fmt.Println("GitHub: https://github.com/layou233/NeverIdle")
}

func parseAndValidateFlags() error {
	flag.Parse()
	return validateFlags()
}

func applyPriority() {
	if *FlagPriority == 666 {
		fmt.Println("[PRIORITY] Use the worst priority by default.")
		controller.SetWorstPriority()
		return
	}

	err := controller.SetPriority(*FlagPriority)
	if err != nil {
		fmt.Println("[PRIORITY] Error when set priority:", err)
	}
}

func startMemoryWaste(nothingEnabled bool) bool {
	if *FlagMemory == 0 {
		return nothingEnabled
	}

	fmt.Println("====================")
	fmt.Println("Starting memory wasting of", *FlagMemory, "GiB")
	go waste.Memory(*FlagMemory)
	runtime.Gosched()
	fmt.Println("====================")
	return false
}

func startCPUWaste(nothingEnabled bool) bool {
	if *FlagCPU != 0 {
		fmt.Println("====================")
		fmt.Println("Starting CPU wasting with interval", *FlagCPU)
		go waste.CPU(*FlagCPU)
		runtime.Gosched()
		fmt.Println("====================")
		return false
	}

	if *FlagCPUPercent == 0 {
		return nothingEnabled
	}

	fmt.Println("====================")
	fmt.Println("Starting CPU wasting with target ratio", *FlagCPUPercent)
	waste.CPUPercent(*FlagCPUPercent)
	runtime.Gosched()
	fmt.Println("====================")
	return false
}

func startNetworkWaste(nothingEnabled bool) bool {
	if *FlagNetwork == 0 {
		return nothingEnabled
	}

	fmt.Println("====================")
	fmt.Println("Starting network speed testing with interval", *FlagNetwork)
	go waste.Network(*FlagNetwork, *FlagNetworkConnectionCount)
	runtime.Gosched()
	fmt.Println("====================")
	return false
}

func waitForShutdown() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	fmt.Println("NeverIdle is running. Press Ctrl+C to stop.")
	<-ctx.Done()
	fmt.Println("Shutting down NeverIdle.")
}

func validateFlags() error {
	if *FlagMemory < 0 {
		return fmt.Errorf("-m must be 0 or greater")
	}
	if *FlagMemory > 0 {
		maxMemoryGiB, err := safeMemoryRequestLimitGiB()
		if err != nil {
			return err
		}
		if *FlagMemory > maxMemoryGiB {
			return fmt.Errorf("-m must be %d GiB or less on this machine", maxMemoryGiB)
		}
	}
	if *FlagCPU < 0 {
		return fmt.Errorf("-c must be 0 or greater")
	}
	if *FlagNetwork < 0 {
		return fmt.Errorf("-n must be 0 or greater")
	}
	if *FlagNetworkConnectionCount <= 0 {
		return fmt.Errorf("-t must be greater than 0")
	}
	if *FlagCPUPercent < 0 || *FlagCPUPercent > 1 {
		return fmt.Errorf("-cp must be a ratio between 0 and 1")
	}
	return nil
}

func safeMemoryRequestLimitGiB() (int, error) {
	vm, err := mem.VirtualMemory()
	if err != nil {
		return 0, fmt.Errorf("failed to inspect system memory: %w", err)
	}

	reserveBytes := uint64(minMemoryReserveGiB) * waste.GiB
	if vm.Available <= reserveBytes {
		return 0, nil
	}

	maxRequestBytes := vm.Available - reserveBytes
	return int(maxRequestBytes / waste.GiB), nil
}
