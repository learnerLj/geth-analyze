// Copyright 2019 The go-ethereum Authors
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

package core

import (
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
)

// txNoncer is a tiny virtual state database to manage the executable nonces of
// accounts in the pool, falling back to reading from a real state database if
// an account is unknown.
type txNoncer struct {
	fallback *state.StateDB
	nonces   map[common.Address]uint64
	//互斥锁 进行加锁操作 保证安全
	lock     sync.Mutex
}

//进行初始化 注意fallback是直接进行的copy 而nonces直接初始化为空map

// newTxNoncer creates a new virtual state database to track the pool nonces.
func newTxNoncer(statedb *state.StateDB) *txNoncer {
	return &txNoncer{
		fallback: statedb.Copy(),
		nonces:   make(map[common.Address]uint64),
	}
}

// get returns the current nonce of an account, falling back to a real state
// database if the account is unknown.
//
func (txn *txNoncer) get(addr common.Address) uint64 {
	// We use mutex for get operation is the underlying
	// state will mutate db even for read access.
	txn.lock.Lock()
	defer txn.lock.Unlock()

	if _, ok := txn.nonces[addr]; !ok {
		txn.nonces[addr] = txn.fallback.GetNonce(addr)
	}
	return txn.nonces[addr]
}

// set inserts a new virtual nonce into the virtual state database to be returned
// whenever the pool requests it instead of reaching into the real state database.
//赋值即可
func (txn *txNoncer) set(addr common.Address, nonce uint64) {
	txn.lock.Lock()
	defer txn.lock.Unlock()

	txn.nonces[addr] = nonce
}

// setIfLower updates a new virtual nonce into the virtual state database if the
// the new one is lower.
//将本地址对应的`nonce`与给的`nonce`进行对比，更新为两者之间最小的那个`nonce`
func (txn *txNoncer) setIfLower(addr common.Address, nonce uint64) {
	txn.lock.Lock()
	defer txn.lock.Unlock()

	if _, ok := txn.nonces[addr]; !ok {
		txn.nonces[addr] = txn.fallback.GetNonce(addr)
	}
	if txn.nonces[addr] <= nonce {
		return
	}
	txn.nonces[addr] = nonce
}

// setAll sets the nonces for all accounts to the given map.
//复制即可
func (txn *txNoncer) setAll(all map[common.Address]uint64) {
	txn.lock.Lock()
	defer txn.lock.Unlock()

	txn.nonces = all
}
