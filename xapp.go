package taskbar

import (
	"encoding/binary"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/xproto"
)

// Reference: https://github.com/linuxmint/xapp/blob/master/libxapp/xapp-gtk-window.c

type xappData struct {
	connection   *xgb.Conn
	progressAtom xproto.Atom
	pulseAtom    xproto.Atom
	window       xproto.Window
}

func xappConnect(xid uint32) (*xappData, error) {
	const progressName = "_NET_WM_XAPP_PROGRESS"
	const pulseName = "_NET_WM_XAPP_PROGRESS_PULSE"

	conn, err := xgb.NewConn()
	if err != nil {
		return nil, err
	}

	data := xappData{}
	data.connection = conn

	progressAtom, err := xproto.InternAtom(conn, false, uint16(len(progressName)), progressName).Reply()
	if err != nil {
		return nil, err
	}
	data.progressAtom = progressAtom.Atom

	pulseAtom, err := xproto.InternAtom(conn, false, uint16(len(pulseName)), pulseName).Reply()
	if err != nil {
		return nil, err
	}
	data.pulseAtom = pulseAtom.Atom

	data.window = xproto.Window(xid)

	return &data, nil
}

func xappDisconnect(data *xappData) error {
	if data == nil {
		return nil
	}
	if err := data.update(0, false); err != nil {
		return err
	}
	data.connection.Close()
	return nil
}

func (data *xappData) update(progress uint64, pulse bool) error {
	var err error
	if progress > 0 {
		err = data.xChangeProperty(data.progressAtom, progress)
	} else {
		err = data.xDeleteProperty(data.progressAtom)
	}
	if err != nil {
		return err
	}

	if pulse {
		err = data.xChangeProperty(data.pulseAtom, 1)
	} else {
		err = data.xDeleteProperty(data.pulseAtom)
	}
	if err != nil {
		return err
	}

	return nil
}

func (data *xappData) xChangeProperty(atom xproto.Atom, value uint64) error {
	bslice := make([]byte, 32)
	binary.LittleEndian.PutUint64(bslice, value)
	err := xproto.ChangePropertyChecked(data.connection, xproto.PropModeReplace, data.window, atom, xproto.AtomCardinal, 32, 8, bslice).Check()
	if err != nil {
		return err
	}
	return nil
}

func (data *xappData) xDeleteProperty(atom xproto.Atom) error {
	if err := xproto.DeletePropertyChecked(data.connection, data.window, atom).Check(); err != nil {
		return err
	}
	return nil
}
