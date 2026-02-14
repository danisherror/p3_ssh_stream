package common

import (
	"encoding/binary"
	"io"
)


const (
        FramePing        uint8 = 1
        FramePong        uint8 = 2
        FrameData        uint8 = 3
        FrameStreamOpen  uint8 = 4
        FrameStreamClose uint8 = 5
)

type Frame struct {
	Type     uint8
	StreamID uint32
	Payload  []byte
}

func EncodeFrame(f Frame) []byte {
	length := uint32(len(f.Payload))

	buf := make([]byte, 1+4+4+len(f.Payload))

	buf[0] = f.Type
	binary.BigEndian.PutUint32(buf[1:5], f.StreamID)
	binary.BigEndian.PutUint32(buf[5:9], length)
	copy(buf[9:], f.Payload)

	return buf
}

func DecodeFrame(r io.Reader) (Frame, error) {
	header := make([]byte, 9)

	_, err := io.ReadFull(r, header)
	if err != nil {
		return Frame{}, err
	}

	frameType := header[0]
	streamID := binary.BigEndian.Uint32(header[1:5])
	length := binary.BigEndian.Uint32(header[5:9])

	payload := make([]byte, length)
	_, err = io.ReadFull(r, payload)
	if err != nil {
		return Frame{}, err
	}

	return Frame{
		Type:     frameType,
		StreamID: streamID,
		Payload:  payload,
	}, nil
}

