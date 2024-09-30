// Package to set window hints like progress in taskbar on Linux.
//
// You can enable progress bar, set an urgent state to an application ("pulse")
// or add a counter badge in taskbar.
// Take a note that you should not make your application rely on these features,
// as whether they will work depends on users desktop environments.
//
// The package uses libunity Launcher API
// (https://wiki.ubuntu.com/Unity/LauncherAPI) by default, which is supported in
// KDE and in GNOME with extensions like Dash-to-Dock. On Cinnamon it uses Xorg
// hints, the same that are used in libxapp (https://github.com/linuxmint/xapp).
// Neither libunity nor libxapp don't need to be installed for this package to
// work, direct Dbus calls and Xorg hints are used instead.
//
// There are 2 environment variables to enforce specific behaviour for this
// library:
//
//   - GO_TASKBAR_BACKEND chooses specific backend no matter what should be
//     selected automatically by code. Only useful for testing.
//     Possible values: libunity, xapp
//   - GO_TASKBAR_DESKTOP_NAME overwrites desktop file name passed to the library
//     for libunity backend to work. This can be useful for packagers (especially
//     for snaps, because desktop file names are changed there automaticaly)
//     without the need to patch a program.
//
// A note for KDE users: while Plasma supports libunity Launcher API, you may
// notice lack of counter badges when developing software. This is caused by
// Plasma behaviour, that expects a counter badge to only be shown by an
// application that also shows notifications. If you application shows them,
// has a proper desktop file and showing badges is allowed in plasma settings,
// counter badge will work as it should.
package taskbar
