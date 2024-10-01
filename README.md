# Go Taskbar package

[![Build workflow](https://github.com/bibelin/taskbar/actions/workflows/go.yml/badge.svg)](https://github.com/bibelin/taskbar/actions/workflows/go.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/bibelin/taskbar)](https://goreportcard.com/report/github.com/bibelin/taskbar) [![Go Reference](https://pkg.go.dev/badge/bibelin/taskbar.svg)](https://pkg.go.dev/bibelin/taskbar)

Set window hints like progress in taskbar on Linux.

```sh
go run example/taskbar-cli.go -desktop <desktop file name> -demo
```

![](screenshots/demo.gif)

Uses libunity Launcher API or Xapp window hints depending on desktop environment. No libs installed needed, Dbus calls and X11 properties are used instead.

![](screenshots/plasma.png)

![](screenshots/mint.png)
