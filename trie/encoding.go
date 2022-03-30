// Copyright 2014 The go-ethereum Authors
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

package trie

// Trie keys are dealt with in three distinct encodings:
//
// KEYBYTES encoding contains the actual key and nothing else. This encoding is the
// input to most API functions.
//
// HEX encoding contains one byte for each nibble of the key and an optional trailing
// 'terminator' byte of value 0x10 which indicates whether or not the node at the key
// contains a value. Hex key encoding is used for nodes loaded in memory because it's
// convenient to access.
//
// COMPACT encoding is defined by the Ethereum Yellow Paper (it's called "hex prefix
// encoding" there) and contains the bytes of the key and a flag. The high nibble of the
// first byte contains the flag; the lowest bit encoding the oddness of the length and
// the second-lowest encoding whether the node at the key is a value node. The low nibble
// of the first byte is zero in the case of an even number of nibbles and the first nibble
// in the case of an odd number. All remaining nibbles (now an even number) fit properly
// into the remaining bytes. Compact encoding is used for nodes stored on disk.

func hexToCompact(hex []byte) []byte {
	terminator := byte(0) //初始化一个值位0的byte
	if hasTerm(hex) {     //验证hex是否有后缀
		terminator = 1         //hex有后缀证明是叶子节点
		hex = hex[:len(hex)-1] //去掉后缀部分的hex编码（16）
	}
	//Compact开辟的空间长度为hex编码的一半再加1，这个1对应的空间是Compact的前缀
	buf := make([]byte, len(hex)/2+1)
	//判断是不是叶子节点
	buf[0] = terminator << 5 // the flag byte
	//这一阶段的buf[0]可以理解为公式中的16*f(t)
	if len(hex)&1 == 1 { //如果长度是奇数
		buf[0] |= 1 << 4 // odd flag//这一阶段的buf[0]可以理解为公式中的16*（f(t)+1）
		buf[0] |= hex[0] // first nibble is contained in the first byte
		hex = hex[1:]
	}
	decodeNibbles(hex, buf[1:])
	return buf
}

// hexToCompactInPlace places the compact key in input buffer, returning the length
// needed for the representation
func hexToCompactInPlace(hex []byte) int {
	var (
		hexLen    = len(hex) // length of the hex input
		firstByte = byte(0)
	)
	// Check if we have a terminator there
	if hexLen > 0 && hex[hexLen-1] == 16 {
		firstByte = 1 << 5
		hexLen-- // last part was the terminator, ignore that
	}
	var (
		binLen = hexLen/2 + 1
		ni     = 0 // index in hex
		bi     = 1 // index in bin (compact)
	)
	if hexLen&1 == 1 {
		firstByte |= 1 << 4 // odd flag
		firstByte |= hex[0] // first nibble is contained in the first byte
		ni++
	}
	for ; ni < hexLen; bi, ni = bi+1, ni+2 {
		hex[bi] = hex[ni]<<4 | hex[ni+1]
	}
	hex[0] = firstByte
	return binLen
}

func compactToHex(compact []byte) []byte {
	if len(compact) == 0 {
		return compact
	}
	//进行展开即可
	base := keybytesToHex(compact)

	// apply terminator flag
	// base[0]包括四种情况
	// 00000000 扩展节点偶数位
	// 00000001 扩展节点奇数位
	// 00000010 叶子节点偶数位
	// 00000011 叶子节点奇数位

	// delete terminator flag
	if base[0] < 2 { //扩展结点
		base = base[:len(base)-1]
	}
	// apply odd flag
	//如果是偶数位，chop等于2，否则等于1
	chop := 2 - base[0]&1
	return base[chop:]
}

func keybytesToHex(str []byte) []byte {
	l := len(str)*2 + 1
	var nibbles = make([]byte, l)
	for i, b := range str {
		nibbles[i*2] = b / 16
		nibbles[i*2+1] = b % 16
	}
	//计入
	nibbles[l-1] = 16
	return nibbles
}

// hexToKeybytes turns hex nibbles into key bytes.
// This can only be used for keys of even length.
func hexToKeybytes(hex []byte) []byte {
	if hasTerm(hex) {
		hex = hex[:len(hex)-1]
	}
	if len(hex)&1 != 0 {
		panic("can't convert hex key of odd length")
	}
	key := make([]byte, len(hex)/2)
	decodeNibbles(hex, key)
	return key
}

func decodeNibbles(nibbles []byte, bytes []byte) {
	for bi, ni := 0, 0; ni < len(nibbles); bi, ni = bi+1, ni+2 {
		bytes[bi] = nibbles[ni]<<4 | nibbles[ni+1]
	}
}

// prefixLen returns the length of the common prefix of a and b.
//返回a&&b相同的前缀长度
func prefixLen(a, b []byte) int {
	var i, length = 0, len(a)
	if len(b) < length {
		length = len(b)
	}
	for ; i < length; i++ {
		if a[i] != b[i] {
			break
		}
	}
	return i
}

// hasTerm returns whether a hex key has the terminator flag.
//判断是否为叶子节点
func hasTerm(s []byte) bool {
	return len(s) > 0 && s[len(s)-1] == 16
}
