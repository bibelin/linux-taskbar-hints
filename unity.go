package taskbar

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
)

// Dbus interface name
const dbusInterface = "com.canonical.Unity.LauncherEntry"

// Dbus introspection scheme
const introspection = `
<node>
	<interface name="` + dbusInterface + `">
		<method name="Query">
			<arg direction="out" type="a{sv}"/>
		</method>
		<signal name="Update">
			<arg name="appUri" type="s"/>
			<arg name="properties" type="a{sv}"/>
		</signal>
	</interface>` + introspect.IntrospectDataString + `</node> `

type libUnityEntry struct {
	connection      *dbus.Conn
	uri             string
	objectPath      dbus.ObjectPath
	progress        float64
	progressVisible bool
	urgent          bool
	count           int64
	countVisible    bool
}

func libUnityConnect(desktopName string) (*libUnityEntry, error) {
	var hash uint64

	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to connect to session bus: %v", err))
	}

	entry := libUnityEntry{}
	if !strings.HasSuffix(desktopName, ".desktop") {
		desktopName = desktopName + ".desktop"
	}
	entry.connection = conn
	entry.uri = "application://" + desktopName

	// DJB hash of URI string is used as identifier
	hash = 5381
	for _, c := range []byte(entry.uri) {
		hash = hash*33 + uint64(c)
	}
	entry.objectPath = dbus.ObjectPath("/com/canonical/unity/launcherentry/" + strconv.FormatUint(hash, 10))

	conn.Export(&entry, entry.objectPath, dbusInterface)
	conn.Export(introspect.Introspectable(introspection), entry.objectPath,
		"org.freedesktop.DBus.Introspectable")

	return &entry, nil
}

func libUnityDisconnect(entry *libUnityEntry) error {
	if entry == nil {
		return nil
	}
	if err := entry.update(0, false, 0); err != nil {
		return err
	}
	if err := entry.connection.Close(); err != nil {
		return err
	}
	return nil
}

// com.canonical.Unity.LauncherEntry.Query Dbus method
func (entry *libUnityEntry) Query() (map[string]interface{}, *dbus.Error) {
	data := map[string]interface{}{
		"progress":         entry.progress,
		"progress-visible": entry.progressVisible,
		"urgent":           entry.urgent,
		"count":            entry.count,
		"countVisible":     entry.countVisible,
	}
	return data, nil
}

func (entry *libUnityEntry) update(progress float64, pulse bool, count int64) error {
	var progressVisible bool
	var countVisible bool

	if progress > 1.0 {
		progress = 1.0
		progressVisible = true
	} else if progress <= 0.0 {
		progress = 0.0
		progressVisible = false
	} else {
		progressVisible = true
	}

	if pulse {
		progressVisible = false
	}

	if count == 0 {
		countVisible = false
	} else {
		countVisible = true
	}

	// Saving properties to use in [Query]
	entry.progress = progress
	entry.progressVisible = progressVisible
	entry.urgent = pulse
	entry.count = count
	entry.countVisible = countVisible

	// Data to send with signal
	data := map[string]interface{}{
		"progress":         progress,
		"progress-visible": progressVisible,
		"urgent":           pulse,
		"count":            count,
		"countVisible":     countVisible,
	}
	// Emit com.canonical.Unity.LauncherEntry.Update signal
	if err := entry.connection.Emit(
		entry.objectPath,
		dbusInterface+".Update",
		entry.uri,
		data); err != nil {
		return err
	}
	return nil
}
