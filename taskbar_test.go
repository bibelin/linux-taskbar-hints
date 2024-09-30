package taskbar

import (
	"os"
	"testing"
)

const desktop = "my.app.desktop"

func TestLibunity(t *testing.T) {
	os.Setenv("GO_TASKBAR_BACKEND", "libunity")
	tb, err := Connect(desktop, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer tb.Disconnect()

	tb.SetProgress(50)
	tb.SetCount(1)
	tb.SetPulse(true)
}

func TestBackendsFail(t *testing.T) {
	os.Setenv("GO_TASKBAR_BACKEND", "libunity")
	tb, err := Connect("", 123)
	if err == nil {
		t.Fail()
	}
	errs := tb.Disconnect()
	if errs == nil {
		t.Fail()
	}

	os.Setenv("GO_TASKBAR_BACKEND", "xapp")
	tb, err = Connect(desktop, 0)
	if err == nil {
		t.Fail()
	}
	errs = tb.Disconnect()
	if errs == nil {
		t.Fail()
	}
}
