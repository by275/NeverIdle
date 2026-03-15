package controller

import (
	"fmt"
	"time"

	"go.einride.tech/pid"
)

const samplingInterval = time.Second

type Device interface {
	Control(value float64)
	Measure() float64
}

func RunPID(
	device Device,
	referenceSignal float64,
	rateImpact float64,
	debug bool,
) *pid.Controller {
	referenceSignal = normalizeReferenceSignal(referenceSignal)
	c := &pid.Controller{
		Config: pid.ControllerConfig{
			ProportionalGain: rateImpact,
			IntegralGain:     1.0,
			DerivativeGain:   0,
		},
	}
	go func() {
		for {
			actualSignal := device.Measure()
			c.Update(pid.ControllerInput{
				ReferenceSignal:  referenceSignal,
				ActualSignal:     actualSignal,
				SamplingInterval: samplingInterval,
			})
			if debug {
				fmt.Printf("actualSignal: %.2f, controlSignal: %.2f\n", actualSignal, c.State.ControlSignal)
			}
			device.Control(c.State.ControlSignal)
		}
	}()
	return c
}

func normalizeReferenceSignal(referenceSignal float64) float64 {
	if referenceSignal < 0 || referenceSignal > 1 {
		fmt.Printf("warning: invalid CPU waste ratio %.2f, falling back to 0.15\n", referenceSignal)
		referenceSignal = 0.15
	}
	return referenceSignal * 100
}
