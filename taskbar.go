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
	xappData   *xappData      // Data for Xapp backend
	progress   int            // progress value (0-100)
	pulse      bool           // whether taskbar pulse is enabled
	count      int            // counter value (only supported by libunity API)
}

// Creates [Taskbar] item.
// `desktopName` is a name of desktop file to be worked with using libunity
// Launcher API (".desktop" suffix can be omitted). `xid` is an xorg window ID
// used in case if taskbar item is modified using xapp window hints, set it to 0
// if not used.
func Connect(desktopName string, xid int) (*Taskbar, error) {
	var t Taskbar
	var session session
	var backend backend

	// Detect session type
	xdgSession := os.Getenv("XDG_SESSION_TYPE")
	switch xdgSession {
	case "wayland":
		session = waylandSession
	case "x11":
		session = xorgSession
	default:
		return nil, errors.New(fmt.Sprintf("Unknown session type: %v", xdgSession))
	}

	// Set backend
	if os.Getenv("GO_TASKBAR_BACKEND") == "libunity" {
		backend = unityBackend
	} else if os.Getenv("GO_TASKBAR_BACKEND") == "xapp" {
		backend = xappBackend
	} else if session == waylandSession {
		backend = unityBackend
	} else if os.Getenv("XDG_CURRENT_DESKTOP") == "X-Cinnamon" {
		backend = xappBackend
	} else {
		backend = unityBackend
	}

	if overrideDesktopName, res := os.LookupEnv("GO_TASKBAR_DESKTOP_NAME"); res {
		desktopName = overrideDesktopName
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
		t = Taskbar{session, backend, entry, nil, 0, false, 0}
	} else {
		xapp, err := xappConnect(uint32(xid))
		if err != nil {
			return nil, err
		}
		t = Taskbar{session, backend, nil, xapp, 0, false, 0}
	}
	return &t, nil
}

// Resets all properties and gracefully disconnects from taskbar
func (t *Taskbar) Disconnect() error {
	if t == nil {
		return errors.New("Not connected to taskbar.")
	}
	if t.backend == xappBackend {
		return xappDisconnect(t.xappData)
	} else {
		return libUnityDisconnect(t.unityEntry)
	}
}

// Gets current progress value
func (t *Taskbar) Progress() int {
	return t.progress
}

// Sets progress value (0-100)
func (t *Taskbar) SetProgress(p int) error {
	if t.progress != p {
		if p > 100 {
			p = 100
		} else if p < 0 {
			p = 0
		}
		t.progress = p
		return t.update()
	}
	return nil
}

// Gets current pulse value
func (t *Taskbar) Pulse() bool {
	return t.pulse
}

// Enables or disables pulse. This property highlights the item in taskbar,
// dragging user attention. If pulse is enabled, progress is not shown.
func (t *Taskbar) SetPulse(p bool) error {
	if t.pulse != p {
		t.pulse = p
		return t.update()
	}
	return nil
}

// Gets current counter value
func (t *Taskbar) Count() int {
	return t.count
}

// Sets counter value (only supported by libunity Launcher API)
func (t *Taskbar) SetCount(c int) error {
	if t.backend == xappBackend {
		return nil
	}
	if t.count != c {
		t.count = c
		return t.update()
	}
	return nil
}

func (t *Taskbar) update() error {
	if t.pulse {
		t.progress = 0
	}
	if t.backend == xappBackend {
		return t.xappData.update(uint64(t.progress), t.pulse)
	} else {
		return t.unityEntry.update(float64(t.progress)/100.0, t.pulse, int64(t.count))
	}
}
