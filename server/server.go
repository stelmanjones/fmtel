package server

import (
	"bytes"
	"encoding/binary"
	"net"

	"github.com/charmbracelet/log"
	"github.com/stelmanjones/fmtel"
)

// Reads telemetry data packets and returns them through provided channel.
func ReadPackets(conn net.PacketConn, ch chan fmtel.ForzaPacket) {
	buf := make([]byte, binary.Size(fmtel.ForzaPacket{}))
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
