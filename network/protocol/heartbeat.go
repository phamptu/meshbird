package protocol

import (
	"fmt"
	"io"
	"net"
)

type (
	HeartbeatMessage []byte
)

func NewHeartbeatMessage(privateIP net.IP) *Packet {
	body := Body{
		Type: TypeHeartbeat,
		Msg:  HeartbeatMessage(privateIP.To4()),
	}
	return &Packet{
		Head: Header{
			Length:  body.Len(),
			Version: CurrentVersion,
		},
		Data: body,
	}
}

func (m HeartbeatMessage) Len() uint16 {
	return uint16(len(m))
}

func (m HeartbeatMessage) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(m)
	return int64(n), err
}

func ReadDecodeHeartbeat(r io.Reader) (HeartbeatMessage, error) {
	logger.Debug("Trying to read Heartbeat message...")

	heartbeatPack, errDecode := ReadAndDecode(r)
	if errDecode != nil {
		logger.WithError(errDecode).Error("Unable to decode package")
		return nil, fmt.Errorf("Error on read Heartbeat package: %v", errDecode)
	}

	if heartbeatPack.Data.Type != TypeHeartbeat {
		return nil, fmt.Errorf("Got non Heartbeat message: %+v", heartbeatPack)
	}

	logger.WithField("msg", heartbeatPack.Data.Msg).Debug("Readed Heartbeat")

	return heartbeatPack.Data.Msg.(HeartbeatMessage), nil
}

func WriteEncodeHeartbeat(w io.Writer, data []byte) (err error) {
	logger.Debug("Trying to write Heartbeat message...")
	if err = EncodeAndWrite(w, NewHeartbeatMessage(data)); err != nil {
		err = fmt.Errorf("Error on write Heartbeat message: %v", err)
	}
	return
}
