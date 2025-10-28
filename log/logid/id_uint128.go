package logid

import (
	"fmt"
	"math/rand/v2"
	"time"
)

const (
	_TIMESTAMP_SHIFT = 16
	_TIMESTAMP_MASK  = 0xFFFFFFFFFFFF

	_SEQUENCE_SHIFT = 6
	_SEQUENCE_MASK  = 0x3FF

	_HIGH_RANDOM_MASK = 0x3F

	base32Table = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"
)

var _ LogID = (*idUint128)(nil)

type idUint128 struct {
	High uint64
	Low  uint64
}

func NewIDUint128(t time.Time, seq uint16) *idUint128 {
	var id idUint128

	timestamp := uint64(t.UnixMilli()) & _TIMESTAMP_MASK
	id.High = timestamp<<_TIMESTAMP_SHIFT | (uint64(seq)&_SEQUENCE_MASK)<<_SEQUENCE_SHIFT

	r := rand.Uint64() & _HIGH_RANDOM_MASK
	id.High |= r

	id.Low = rand.Uint64()
	return &id
}

func (id *idUint128) String() string {
	buf := make([]byte, 0, 34)
	timestamp := id.High >> _TIMESTAMP_SHIFT

	mill := timestamp % 1000
	t := time.UnixMilli(int64(timestamp))
	timeStr := t.Format("20060102150405") + fmt.Sprintf("%03d", mill)

	seq := (id.High >> _SEQUENCE_SHIFT) & _SEQUENCE_MASK
	seqStr := fmt.Sprintf("%03X", uint64(seq))

	rands := encode70Bits(id.High, id.Low, _HIGH_RANDOM_MASK)

	buf = append(buf, []byte(timeStr)...)
	buf = append(buf, []byte(seqStr)...)
	buf = append(buf, rands...)
	return string(buf)
}

func encode70Bits(high, low uint64, highMask uint64) []byte {
	// 提取high的低6位有效信息
	highPart := high & highMask // highMask=0x3F（6位）
	result := make([]byte, 14)  // 70位=14×5位，结果为14位

	// 遍历14个5位组（i=0到13）
	for i := range 14 {
		// 计算当前5位组在70位中的最高位和最低位位置（70位范围：69（最高）~0（最低））
		highBit := 69 - i*5   // 第i组的最高位
		lowBit := highBit - 4 // 第i组的最低位

		var val uint64

		if highBit >= 64 {
			// 部分或全部在highPart中（highPart占70位的64~69位，共6位）
			// highPart的内部位索引：69→5，68→4，...，64→0（H5~H0）
			hpHigh := highBit - 64 // 转换为highPart内部的高位索引（0~5）
			hpLow := lowBit - 64   // 转换为highPart内部的低位索引

			if hpLow < 0 {
				// 跨highPart和low：部分在highPart（hpHigh~0），部分在low（63~lowBit）
				// 1. 提取highPart中hpHigh~0的比特（共hpHigh+1位）
				highVal := (highPart >> 0) & ((1 << (hpHigh + 1)) - 1)
				// 2. 提取low中63~lowBit的比特（共63 - lowBit + 1位）
				lowBits := 63 - lowBit + 1
				lowVal := (low >> lowBit) & ((1 << lowBits) - 1)
				// 合并：highVal左移lowBits位，再或上lowVal
				val = (highVal << lowBits) | lowVal
			} else {
				// 完全在highPart中：提取hpHigh~hpLow的比特（共5位）
				val = (highPart >> hpLow) & ((1 << (hpHigh - hpLow + 1)) - 1)
			}
		} else {
			// 完全在low中：提取highBit~lowBit的5位（直接右移lowBit后取低5位）
			val = (low >> lowBit) & 0x1F // 0x1F是5位掩码（11111）
		}

		// 映射到Base32字符
		result[i] = base32Table[val]
	}

	return result
}
