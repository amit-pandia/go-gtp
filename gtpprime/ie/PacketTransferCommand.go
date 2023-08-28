package ie

import (
	"io"
)

// NewPacketTransferCommand creates a new Packet Transfer Command IE.
func NewPacketTransferCommand(packetTransferCommand uint8) *IE {
	return newUint8ValIE(PacketTransferCommand, packetTransferCommand)
}

// PacketTransferCommand returns PacketTransferCommand value if the type of IE matches.
func (i *IE) PacketTransferCommand() (uint8, error) {
	if i.Type != PacketTransferCommand {
		return 0, &InvalidTypeError{Type: i.Type}
	}
	if len(i.Payload) < 1 {
		return 0, io.ErrUnexpectedEOF
	}
	return i.Payload[0], nil
}

// MustPacketTransferCommand returns PacketTransferCommand in uint8, ignoring errors.
// This should only be used if it is assured to have the value.
func (i *IE) MustPacketTransferCommand() uint8 {
	v, _ := i.PacketTransferCommand()
	return v
}
