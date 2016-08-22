package timer

import (
	"testing"
	"time"
)

func TestLeftTime(t *testing.T) {
	initialTimeout := 10
	timeNanoseconds := time.Second * time.Duration(initialTimeout)
	tim := InitHourglass(timeNanoseconds)

	<-time.After(time.Second * time.Duration(1))
	if tim.GetLeftNanoseconds() < time.Second*time.Duration(8) ||
		tim.GetLeftNanoseconds() > time.Second*time.Duration(9) {
		t.Errorf("Testing timer.GetLeftTime for 1 sec err, TimeLeft: %d", int64(tim.GetLeftNanoseconds()))
	}

	<-time.After(time.Second * time.Duration(3))
	if tim.GetLeftNanoseconds() < time.Second*time.Duration(6) ||
		tim.GetLeftNanoseconds() > time.Second*time.Duration(7) {
		t.Errorf("Testing timer.GetLeftTime for 3 sec err, TimeLeft: %d", int64(tim.GetLeftNanoseconds()))
	}
}

func TestOverflow(t *testing.T) {
	initialTimeout := 10
	tim := InitHourglass(time.Second * time.Duration(initialTimeout))

	<-time.After(time.Second * time.Duration(11))
	if int(tim.GetLeftNanoseconds()) > 0 {
		t.Errorf("Testing timer.GetLeftTime overflow err, TimeLeft: %d", int64(tim.GetLeftNanoseconds()))
	}
}
