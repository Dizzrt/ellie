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

	localIPHexStr = make([]byte, 0, 8)
	for _, b := range ip {
		localIPHexStr = append(localIPHexStr, fmt.Sprintf("%02X", b)...)
	}
}

var _ LogID = (*ID128Bits)(nil)
var (
	length        = 36
	localIPHexStr []byte
)

type ID128Bits []byte

func NewID128Bits(t time.Time, sequence int) LogID {
	// 0~47bits mill-timestamp 17
	// 48~51bits reserved 1
	// 52~63bits sequence 3
	// 64~95bits local ip 8
	// 96~127bits random

	buf := make([]byte, length)

	copy(buf[0:], t.AppendFormat(nil, "20060102150405"))
	copy(buf[14:], fmt.Appendf(nil, "%03d", t.UnixMilli()%1000))
	copy(buf[17:], fmt.Appendf(nil, "0%03X", sequence))
	copy(buf[21:], localIPHexStr)

	var temp [4]byte
	binary.BigEndian.PutUint32(temp[:], rand.Uint32())
	base32.StdEncoding.WithPadding(base32.NoPadding).Encode(buf[29:], temp[:])

	return ID128Bits(buf)
}

func (id ID128Bits) String() string {
	return string(id)
}
