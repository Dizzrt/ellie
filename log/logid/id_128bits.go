package logid

import (
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"math/rand/v2"
	"net"
	"time"
)

func init() {
	dial, err := net.Dial("udp", "8.8.8.8:8")
	if err != nil {
		panic(err)
	}

	ip := dial.LocalAddr().(*net.UDPAddr).IP.To4()

	localIP = make([]byte, 0, 8)
	for _, b := range ip {
		localIP = append(localIP, fmt.Sprintf("%02X", b)...)
	}
}

const (
	_FLAGS_SHIFT   = 12
	_FLAGS_MASK    = 0xF
	_SEQUENCE_MASK = 0xFFF

	flags = (0xF & _FLAGS_MASK) << _FLAGS_SHIFT // 0xF is reserved for future use
)

var _ LogID = (*ID128Bits)(nil)
var (
	length  = 36
	localIP []byte
)

type ID128Bits []byte

func NewID128Bits(t time.Time, sequence int) LogID {
	// 0~47bits mill-timestamp, 17 chars
	// 48~51bits reserved(flags), 1 chars
	// 52~63bits sequence, 3 chars
	// 64~95bits local ip, 8 chars
	// 96~127bits random, encode by base32, 7 chars
	// total 128 bits, 36 chars

	buf := make([]byte, length)

	copy(buf[0:], t.AppendFormat(nil, "20060102150405"))
	copy(buf[14:], fmt.Appendf(nil, "%03d%04X", t.UnixMilli()%1000, flags|(sequence&_SEQUENCE_MASK)))
	copy(buf[21:], localIP)

	var temp [4]byte
	binary.BigEndian.PutUint32(temp[:], rand.Uint32())
	base32.StdEncoding.WithPadding(base32.NoPadding).Encode(buf[29:], temp[:])

	return ID128Bits(buf)
}

func (id ID128Bits) String() string {
	return string(id)
}
