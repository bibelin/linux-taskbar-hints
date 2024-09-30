// Package to set window hints like progress in taskbar on Linux.
package taskbar

import (
	"errors"
	"fmt"
	"os"
)

type session int
type backend int

const (
	waylandSession session = iota
	xorgSession
)
const (
	unityBackend backend = iota
	xappBackend
)

// Taskbar item data
type Taskbar struct {
	session    session        // session type (wayland or xorg)
	backend    backend        // taskbar backend (libunity or xapp)
	unityEntry *libUnityEntry // LibUnity launcher entry
	xid        int32          // xorg window ID
	progress   int            // progress value (0-100)
	pulse      bool           // whether taskbar pulse is enabled
	count      int            // counter value (only supported by libunity API)
}

// Creates [Taskbar] item.
// `desktopName` is a name of desktop file to be worked with using libunity
// Launcher API (".desktop" suffix can be omitted). `xid` is an xorg window ID
// used in case if taskbar item is modified using xapp window hints.
func Connect(desktopName string, xid int32) (*Taskbar, error) {
	var t Taskbar
	var session session
	var backend backend

	// Detect session type
	switch os.Getenv("XDG_SESSION_TYPE") {
	case "wayland":
		session = waylandSession
	case "x11":
		session = xorgSession
	_:
		return nil, errors.New(fmt.Sprintf("Unknown session type: %v", session))
	}

	// Set backend
	if os.Getenv("GO_TASKBAR_BACKEND") == "libunity" {
		backend = unityBackend
	} else if os.Getenv("GO_TASKBAR_BACKEND") == "xapp" {
		backend = xappBackend
	} else if session == waylandSession {
		backend = unityBackend
	} else if os.Getenv("XDG_CURRENT_DESKTOP") == "Cinnamon" {
		backend = xappBackend
	} else {
		backend = unityBackend
	}

	// Check if current backend can be used
	if backend == unityBackend && desktopName == "" {
		return nil, errors.New("LibUnity backend was chosen, but desktop file name is empty.")
	}
	if backend == xappBackend && xid == 0 {
		return nil, errors.New("Xapp backend was chosen, but XID isn't provided.")
	}

	if backend == unityBackend {
		entry, err := libUnityConnect(desktopName)
		if err != nil {
			return nil, err
		}
		t = Taskbar{session, backend, entry, xid, 0, false, 0}
	} else {
		t = Taskbar{session, backend, nil, xid, 0, false, 0}
	}
	return &t, nil
}

func (t *Taskbar) Disconnect() []error {
	if t == nil {
		return []error{errors.New("Not connected to taskbar.")}
	}
	if t.backend == xappBackend {
		// TODO: Xapp implementation
		return nil
	} else {
		return libUnityDisconnect(t.unityEntry)
	}
}

// Get current progress value
func (t *Taskbar) Progress() int {
	return t.progress
}

// Set progress value (0-100)
func (t *Taskbar) SetProgress(p int) error {
	t.progress = p
	return t.update()
}

// Get current pulse value
func (t *Taskbar) Pulse() bool {
	return t.pulse
}

// Enable or disable pulse. This mode "highlights" the item in taskbar, dragging
// user attention. If pulse is enabled, progress is not shown.
func (t *Taskbar) SetPulse(p bool) error {
	t.pulse = p
	return t.update()
}

// Get current counter value
func (t *Taskbar) Count() int {
	return t.count
}

// Set counter value (only supported by libunity Launcher API)
func (t *Taskbar) SetCount(c int) error {
	t.count = c
	return t.update()
}

func (t *Taskbar) update() error {
	if t.backend == xappBackend {
		// TODO: Xapp implementation
		return nil
	} else {
		return t.unityEntry.update(float64(t.progress)/100.0, t.pulse, int64(t.count))
	}
}
