// Package to set window hints like progress in taskbar on Linux.
// It uses libunity Launcher API (https://wiki.ubuntu.com/Unity/LauncherAPI)
// by default, which is supported in KDE and GNOME with extensions like
// Dash-to-Dock. On Cinnamon it uses Xorg hints, the same that are used in
// libxapp (https://projects.linuxmint.com/xapp/reference/).
// Neither libunity nor libxapp don't need to be installed for this package to
// work, direct Dbus calls and Xorg hints are used instead.
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
	count      uint64         // counter value (only supported by libunity API)
}

// Creates [Taskbar] item, returns pointer to it and any error if happened.
// `desktopName` is a name of desktop file to be worked with using libunity
// Launcher API (".desktop" suffix can be omitted). `xid` is an xorg window ID
// used in case if taskbar item is modified using xapp window hints.
func Init(desktopName string, xid int32) (*Taskbar, error) {
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
	if session == waylandSession {
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
		entry, _ := libUnityInit(desktopName)
		t = Taskbar{session, backend, entry, xid, 0, false, 0}
	} else {
		t = Taskbar{session, backend, nil, xid, 0, false, 0}
	}
	return &t, nil
}

// Get current progress value
func (t *Taskbar) Progress() int {
	return t.progress
}

// Set progress value (0-100)
func (t *Taskbar) SetProgress(p int) {
	t.progress = p
	t.update()
}

// Get current pulse value
func (t *Taskbar) Pulse() bool {
	return t.pulse
}

// Enable or disable pulse. This mode "highlights" the item in taskbar, dragging
// user attention. If pulse is enabled, progress is not shown.
func (t *Taskbar) SetPulse(p bool) {
	t.pulse = p
	t.update()
}

// Get current counter value
func (t *Taskbar) Count() uint64 {
	return t.count
}

// Set counter value (only supported by libunity Launcher API)
func (t *Taskbar) SetCount(c uint64) {
	t.count = c
	t.update()
}

func (t *Taskbar) update() error {
	if t.backend == xappBackend {
		// TODO: Xapp implementation
		return nil
	} else {
		return t.unityEntry.update(float64(t.progress)/100.0, t.pulse, t.count)
	}
}
