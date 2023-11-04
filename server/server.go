package server

import (
	"bytes"
	"encoding/binary"
	"net"
	"time"

	"github.com/charmbracelet/log"
	"github.com/stelmanjones/fmtel"
)

// Reads telemetry data packets and returns them through provided channel.
// BUG: RefreshRate makes rendering fall behind. Place ticker in event loop.
func ReadPackets(conn net.PacketConn, ch chan fmtel.ForzaPacket, refreshRate int) {
	buf := make([]byte, binary.Size(fmtel.ForzaPacket{}))
	delay := 1000 / refreshRate
	ticker := time.NewTicker(time.Duration(delay * int(time.Millisecond)))
	defer ticker.Stop()
	for {
		_, _, err := conn.ReadFrom(buf)
		if err != nil {
			log.Error(err)
		}
		var packet fmtel.ForzaPacket
		err = binary.Read(bytes.NewReader(buf), binary.LittleEndian, &packet)
		if err != nil {
			log.Error(err)
		}
		ch <- packet

	}
}

// HACK: Remove ASAP.
func DummyUDP(ch chan fmtel.ForzaPacket) {
	data := fmtel.DefaultForzaPacket
	data.IsRaceOn = 1

	ticker := time.NewTicker(time.Duration(100 * time.Millisecond))
	defer ticker.Stop()
	for {
		<-ticker.C
		data.CurrentEngineRpm += 1.0
		if data.CurrentEngineRpm >= 9999.0 {
			data.CurrentEngineRpm = 0.0
		}
		ch <- data

	}
}
