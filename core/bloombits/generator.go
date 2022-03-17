// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package bloombits

import (
	"errors"

	"github.com/ethereum/go-ethereum/core/types"
)

var (
	// errSectionOutOfBounds is returned if the user tried to add more bloom filters
	// to the batch than available space, or if tries to retrieve above the capacity.
	errSectionOutOfBounds = errors.New("section out of bounds") //段越界或者查询的数据过多

	// errBloomBitOutOfBounds is returned if the user tried to retrieve specified
	// bit bloom above the capacity.
	errBloomBitOutOfBounds = errors.New("bloom bit out of bounds") //位索引越界
)

/*因为每个区块的 bloom 有 2048 位，其中每个元素哈希后只有三个位置置为1，如果逐个区块检索效率较低。
因此将原本多个布隆过滤器的 bit 矩阵：
[A0, A1, ..., A2047]
[B0, B1, ..., B2047]
...
[H0, H1, ..., H2047]
转置为：
[A0, B0, ..., H0]
[A1, B1, ..., H1]
...
[A2047, B2047, ..., H2047]
其中 blooms 远不止 A-G 八个布隆过滤器，而是共 4096 个计数器，恰好对应这一个 section 的区块（4096个）
*/

// Generator takes a number of bloom filters and generates the rotated bloom bits
// to be used for batched filtering.
type Generator struct {
	//转置后的 bloom 数据，实际上为 2048*(sections/8) 的矩阵
	blooms [types.BloomBitLength][]byte // Rotated blooms for per-bit matching
	//段的个数，也是布隆过滤器的个数
	sections uint // Number of sections to batch together
	//当前批量处理中的下一个将处理的段，也就是下一个 bloom
	nextSec uint // Next section to set when adding a bloom
}

// NewGenerator creates a rotated bloom generator that can iteratively fill a
// batched bloom filter's bits.
func NewGenerator(sections uint) (*Generator, error) {
	//段的数量需要是 8 的倍数，这样可以恰好按比特填充到字节数组里。
	if sections%8 != 0 {
		return nil, errors.New("section count not multiple of 8")
	}
	b := &Generator{sections: sections}
	//请注意转置矩阵。这里因为一个字节占 8 位，而这里使用 byte 数组，因此只要占 1/8 的位置
	for i := 0; i < types.BloomBitLength; i++ {
		b.blooms[i] = make([]byte, sections/8)
	}
	return b, nil
}

// AddBloom takes a single bloom filter and sets the corresponding bit column
// in memory accordingly.
func (b *Generator) AddBloom(index uint, bloom types.Bloom) error {
	// Make sure we're not adding more bloom filters than our capacity
	if b.nextSec >= b.sections {
		//超过了设定的批量处理的段的个数
		return errSectionOutOfBounds
	}
	//index 与下个将处理的段对应
	if b.nextSec != index {
		return errors.New("bloom filter with unexpected index")
	}

	// Rotate the bloom and insert into our collection
	byteIndex := b.nextSec / 8 //根据下一个 section 的编号找到字节数组中的索引
	//在一个字节中的比特的索引，注意一个字节中存有 8 sections 的
	// bloom 的比特向量的某一位
	bitIndex := byte(7 - b.nextSec%8)
	//开始初始化 每一列对应的 bloom
	for byt := 0; byt < types.BloomByteLength; byt++ {
		//八位一组，大端序
		bloomByte := bloom[types.BloomByteLength-1-byt]
		if bloomByte == 0 {
			continue
		}
		base := 8 * byt
		b.blooms[base+7][byteIndex] |= ((bloomByte >> 7) & 1) << bitIndex
		b.blooms[base+6][byteIndex] |= ((bloomByte >> 6) & 1) << bitIndex
		b.blooms[base+5][byteIndex] |= ((bloomByte >> 5) & 1) << bitIndex
		b.blooms[base+4][byteIndex] |= ((bloomByte >> 4) & 1) << bitIndex
		b.blooms[base+3][byteIndex] |= ((bloomByte >> 3) & 1) << bitIndex
		b.blooms[base+2][byteIndex] |= ((bloomByte >> 2) & 1) << bitIndex
		b.blooms[base+1][byteIndex] |= ((bloomByte >> 1) & 1) << bitIndex
		b.blooms[base][byteIndex] |= (bloomByte & 1) << bitIndex
	}
	b.nextSec++
	return nil
}

//返回 blooms 的某一行，注意这不是一个布隆过滤器，而是连续的若干个 section 的 bloom 的索引为 idx 构成的集合

// Bitset returns the bit vector belonging to the given bit index after all
// blooms have been added.
func (b *Generator) Bitset(idx uint) ([]byte, error) {
	//因为 nextSec 递增，可以表示是否完成了给 Generator 的赋值
	if b.nextSec != b.sections {
		return nil, errors.New("bloom not fully generated yet")
	}
	if idx >= types.BloomBitLength {
		return nil, errBloomBitOutOfBounds
	}
	return b.blooms[idx], nil
}
