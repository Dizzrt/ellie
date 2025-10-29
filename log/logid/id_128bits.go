package logid

import (
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"math/rand/v2"
	"net"
	"sync"
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

	hexTable = "0123456789ABCDEF"
	flags    = (0xF & _FLAGS_MASK) << _FLAGS_SHIFT // 0xF is reserved for future use
)

var _ LogID = (*ID128Bits)(nil)
var (
	length  = 36
	localIP []byte

	bufPool = sync.Pool{
		New: func() any {
			buf := make([]byte, length)
			return &buf
		},
	}

	randPool = sync.Pool{
		New: func() any {
			r := rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))
			return r
		},
	}
)

type ID128Bits []byte

func writeTimeDigits(bufPtr *[]byte, offset int, value int) {
	tens := value / 10
	ones := value % 10
	(*bufPtr)[offset] = byte(tens + '0')
	(*bufPtr)[offset+1] = byte(ones + '0')
}

func NewID128Bits(t time.Time, sequence int) LogID {
	// 0~47bits mill-timestamp, 17 chars
	// 48~51bits reserved(flags), 1 chars
	// 52~63bits sequence, range [0, 0xFFF], 3 chars
	// 64~95bits local ip (only ipv4), 8 chars
	// 96~127bits random, encode by base32, 7 chars
	// total 128 bits, 36 chars

	bufPtr := bufPool.Get().(*[]byte)
	defer bufPool.Put(bufPtr)

	year := t.Year()
	(*bufPtr)[0] = byte(year/1000 + '0')
	(*bufPtr)[1] = byte((year/100)%10 + '0')
	(*bufPtr)[2] = byte((year/10)%10 + '0')
	(*bufPtr)[3] = byte(year%10 + '0')

	writeTimeDigits(bufPtr, 4, int(t.Month()))
	writeTimeDigits(bufPtr, 6, t.Day())
	writeTimeDigits(bufPtr, 8, t.Hour())
	writeTimeDigits(bufPtr, 10, t.Minute())
	writeTimeDigits(bufPtr, 12, t.Second())

	ms := t.UnixMilli() % 1000
	(*bufPtr)[14] = byte(ms/100 + '0')
	(*bufPtr)[15] = byte((ms/10)%10 + '0')
	(*bufPtr)[16] = byte(ms%10 + '0')

	// 4096 requests per millisecond in a single node is enough for 99.9999...% cases
	// so we dont need to handle the case when seq exceeds 0xFFF, just mask it
	// besides, there is a random number at the end
	seqHex := flags | (sequence & _SEQUENCE_MASK)
	(*bufPtr)[17] = hexTable[(seqHex>>12)&0xF]
	(*bufPtr)[18] = hexTable[(seqHex>>8)&0xF]
	(*bufPtr)[19] = hexTable[(seqHex>>4)&0xF]
	(*bufPtr)[20] = hexTable[seqHex&0xF]

	copy((*bufPtr)[21:], localIP)

	r := randPool.Get().(*rand.Rand)
	defer randPool.Put(r)

	var temp [4]byte
	binary.BigEndian.PutUint32(temp[:], r.Uint32())
	base32.StdEncoding.WithPadding(base32.NoPadding).Encode((*bufPtr)[29:], temp[:])

	return ID128Bits((*bufPtr)[:length])
}

func (id ID128Bits) String() string {
	return string(id)
}
