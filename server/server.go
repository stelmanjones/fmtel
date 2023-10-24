package server

import (
	"bytes"
	"encoding/binary"
	"net"

	"github.com/charmbracelet/log"
	"github.com/stelmanjones/fmtel"
)

func ReadPackets(conn net.PacketConn, ch chan pkg.ForzaPacket) {
	buf := make([]byte, binary.Size(pkg.ForzaPacket{}))
	for {
		_, _, err := conn.ReadFrom(buf)
		if err != nil {
			log.Error(err)
		}
		var packet pkg.ForzaPacket
		err = binary.Read(bytes.NewReader(buf), binary.LittleEndian, &packet)
		if err != nil {
			log.Error(err)
		}
		ch <- packet

	}
}
