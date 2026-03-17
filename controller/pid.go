package controller

import (
	"time"

	"github.com/by275/neveridle/internal/log"
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
				log.Logf("PID", "actualSignal=%.2f controlSignal=%.2f", actualSignal, c.State.ControlSignal)
			}
			device.Control(c.State.ControlSignal)
		}
	}()
	return c
}

func normalizeReferenceSignal(referenceSignal float64) float64 {
	if referenceSignal < 0 || referenceSignal > 1 {
		log.Logf("PID", "invalid CPU waste ratio %.2f, falling back to 0.15", referenceSignal)
		referenceSignal = 0.15
	}
	return referenceSignal * 100
}
