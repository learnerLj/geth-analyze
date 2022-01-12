// Package checkpointoracle is an on-chain light client checkpoint oracle.

package checkpointoracle

//go:generate abigen --sol contract/oracle.sol --pkg contract --out contract/oracle.go

//使用 abigen 工具根据目录contract 下的 oracle.sol,在 contract 包内 生成目录contract 下的 oracle.go 文件，
//里面是合约相关的 Golang语言的封装

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/checkpointoracle/contract"
	"github.com/ethereum/go-ethereum/core/types"
)

// CheckpointOracle is a Go wrapper around an on-chain checkpoint oracle contract.
type CheckpointOracle struct {
	address  common.Address
	contract *contract.CheckpointOracle //预言机类型封装，包括了调用内容、绑定的合约的封装、筛选器
}

//绑定作为检查点的合约，返回封装好的合约实例

// NewCheckpointOracle binds checkpoint contract and returns a registrar instance.
func NewCheckpointOracle(contractAddr common.Address, backend bind.ContractBackend) (*CheckpointOracle, error) {
	c, err := contract.NewCheckpointOracle(contractAddr, backend)
	if err != nil {
		return nil, err
	}
	return &CheckpointOracle{address: contractAddr, contract: c}, nil
}

//获取地址

// ContractAddr returns the address of contract.
func (oracle *CheckpointOracle) ContractAddr() common.Address {
	return oracle.address
}

//获取可直接用于调用函数的合约实例

// Contract returns the underlying contract instance.
func (oracle *CheckpointOracle) Contract() *contract.CheckpointOracle {
	return oracle.contract
}

//查找某一段内生成检查点时的投票事件（即参与验证签名）

// LookupCheckpointEvents searches checkpoint event for specific section in the
// given log batches.
func (oracle *CheckpointOracle) LookupCheckpointEvents(blockLogs [][]*types.Log, section uint64, hash common.Hash) []*contract.CheckpointOracleNewCheckpointVote {
	var votes []*contract.CheckpointOracleNewCheckpointVote

	for _, logs := range blockLogs { //需检索的日志
		for _, log := range logs {
			event, err := oracle.contract.ParseNewCheckpointVote(*log) //解析日志中的事件
			if err != nil {
				continue
			}
			if event.Index == section && event.CheckpointHash == hash { //事件在需要检索的这一段，并且哈希值正确。
				votes = append(votes, event)
			}
		}
	}
	return votes
}

//创建检查点，创建时获取发起的签名，然后调用根据 oracle 合约生成的封装好的代码，给合约发消息

// RegisterCheckpoint registers the checkpoint with a batch of associated signatures
// that are collected off-chain and sorted by lexicographical order.
//
// Notably all signatures given should be transformed to "ethereum style" which transforms
// v from 0/1 to 27/28 according to the yellow paper.
func (oracle *CheckpointOracle) RegisterCheckpoint(opts *bind.TransactOpts, index uint64, hash []byte, rnum *big.Int, rhash [32]byte, sigs [][]byte) (*types.Transaction, error) {
	var (
		r [][32]byte
		s [][32]byte
		v []uint8
	)
	for i := 0; i < len(sigs); i++ {
		if len(sigs[i]) != 65 {
			return nil, errors.New("invalid signature")
		}
		r = append(r, common.BytesToHash(sigs[i][:32]))
		s = append(s, common.BytesToHash(sigs[i][32:64]))
		v = append(v, sigs[i][64])
	}
	return oracle.contract.SetCheckpoint(opts, rnum, rhash, common.BytesToHash(hash), index, v, r, s)
}
