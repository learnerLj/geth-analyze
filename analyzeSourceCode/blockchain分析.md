>  参考： [以太坊blockchain源码分析 - mindcarver - 博客园 (cnblogs.com)](https://www.cnblogs.com/1314xf/articles/13827186.html)

# [以太坊blockchain源码分析](https://www.cnblogs.com/1314xf/p/13827186.html)

## blockchain关键元素

- db：持久化到底层数据储存，即leveldb；
- genesisBlock：创始区块
- currentBlock：当前区块，blockchain中并不是储存链所有的block，而是通过currentBlock向前回溯直到genesisBlock，这样就构成了区块链
- bodyCache、bodyRLPCache、blockCache、futureBlocks：区块链中的缓存结构，用于加快区块链的读取和构建；
- hc：headerchain区块头链，由blockchain额外维护的另一条链，由于Header和Block的储存空间是有很大差别的，但同时Block的Hash值就是Header（RLP）的Hash值，所以维护一个headerchain可以用于快速延长链，验证通过后再下载blockchain，或者可以与blockchain进行相互验证；
- processor：执行区块链交易的接口，收到一个新的区块时，要对区块中的所有交易执行一遍，一方面是验证，一方面是更新世界状态；
- validator：验证数据有效性的接口
- futureBlocks：收到的区块时间大于当前头区块时间15s而小于30s的区块，可作为当前节点待处理的区块。

------

## 函数介绍

```go
// BadBlocks 处理客户端从网络上获取的最近的bad block列表
func (bc *BlockChain) BadBlocks() []*types.Block {}

// addBadBlock 把bad block放入缓存
func (bc *BlockChain) addBadBlock(block *types.Block) {}
// CurrentBlock取回主链的当前头区块，这个区块是从blockchian的内部缓存中取得
func (bc *BlockChain) CurrentBlock() *types.Block {}
 
// CurrentHeader检索规范链的当前头区块header。从HeaderChain的内部缓存中检索标头。
func (bc *BlockChain) CurrentHeader() *types.Header{}
 
// CurrentFastBlock取回主链的当前fast-sync头区块，这个区块是从blockchian的内部缓存中取得
func (bc *BlockChain) CurrentFastBlock() *types.Block {}
// 将活动链或其子集写入给定的编写器.
func (bc *BlockChain) Export(w io.Writer) error {}
func (bc *BlockChain) ExportN(w io.Writer, first uint64, last uint64) error {}
// FastSyncCommitHead快速同步，将当前头块设置为特定hash的区块。
func (bc *BlockChain) FastSyncCommitHead(hash common.Hash) error {}
// GasLimit返回当前头区块的gas limit
func (bc *BlockChain) GasLimit() uint64 {}
// Genesis 取回genesis区块
func (bc *BlockChain) Genesis() *types.Block {}
// 通过hash从数据库或缓存中取到一个区块体(transactions and uncles)或RLP数据
func (bc *BlockChain) GetBody(hash common.Hash) *types.Body {}
func (bc *BlockChain) GetBodyRLP(hash common.Hash) rlp.RawValue {}
// GetBlock 通过hash和number取到区块
func (bc *BlockChain) GetBlock(hash common.Hash, number uint64) *types.Block {}
// GetBlockByHash 通过hash取到区块
func (bc *BlockChain) GetBlockByHash(hash common.Hash) *types.Block {}
// GetBlockByNumber 通过number取到区块
func (bc *BlockChain) GetBlockByNumber(number uint64) *types.Block {}
// 获取给定hash和number区块的header
func (bc *BlockChain) GetHeader(hash common.Hash, number uint64) *types.Header{}
 
// 获取给定hash的区块header
func (bc *BlockChain) GetHeaderByHash(hash common.Hash) *types.Header{}
 
// 获取给定number的区块header
func (bc *BlockChain) GetHeaderByNumber(number uint64) *types.Header{}
// HasBlock检验hash对应的区块是否完全存在数据库中
func (bc *BlockChain) HasBlock(hash common.Hash, number uint64) bool {}
 
// 检查给定hash和number的区块的区块头是否存在数据库
func (bc *BlockChain) HasHeader(hash common.Hash, number uint64) bool{}
 
// HasState检验state trie是否完全存在数据库中
func (bc *BlockChain) HasState(hash common.Hash) bool {}
 
// HasBlockAndState检验hash对应的block和state trie是否完全存在数据库中
func (bc *BlockChain) HasBlockAndState(hash common.Hash, number uint64) bool {}
// 获取给定hash的区块的总难度
func (bc *BlockChain) GetTd(hash common.Hash, number uint64) *big.Int{}
// 获取从给定hash的区块到genesis区块的所有hash
func (bc *BlockChain) GetBlockHashesFromHash(hash common.Hash, max uint64) []common.Hash{}
 
// GetReceiptsByHash 在特定的区块中取到所有交易的收据
func (bc *BlockChain) GetReceiptsByHash(hash common.Hash) types.Receipts {}
 
// GetBlocksFromHash 取到特定hash的区块及其n-1个父区块
func (bc *BlockChain) GetBlocksFromHash(hash common.Hash, n int) (blocks []*types.Block) {}
 
// GetUnclesInChain 取回从给定区块到向前回溯特定距离到区块上的所有叔区块
func (bc *BlockChain) GetUnclesInChain(block *types.Block, length int) []*types.Header {}
// insert 将新的头块注入当前块链。 该方法假设该块确实是真正的头。
// 如果它们较旧或者它们位于不同的侧链上，它还会将头部标题和头部快速同步块重置为同一个块。
func (bc *BlockChain) insert(block *types.Block) {}
 
// InsertChain尝试将给定批量的block插入到规范链中，否则，创建一个分叉。 如果返回错误，它将返回失败块的索引号以及描述错误的错误。
//插入完成后，将触发所有累积的事件。
func (bc *BlockChain) InsertChain(chain types.Blocks) (int, error){}
 
// insertChain将执行实际的链插入和事件聚合。 
// 此方法作为单独方法存在的唯一原因是使用延迟语句使锁定更清晰。
func (bc *BlockChain) insertChain(chain types.Blocks) (int, []interface{}, []*types.Log, error){}
 
// InsertHeaderChain尝试将给定的headerchain插入到本地链中，可能会创建一个重组
func (bc *BlockChain) InsertHeaderChain(chain []*types.Header, checkFreq int) (int, error){}
 
// InsertReceiptChain 使用交易和收据数据来完成已经存在的headerchain
func (bc *BlockChain) InsertReceiptChain(blockChain types.Blocks, receiptChain []types.Receipts) (int, error) {}
//loadLastState从数据库加载最后一个已知的链状态。
func (bc *BlockChain) loadLastState() error {}
// Processor 返回当前current processor.
func (bc *BlockChain) Processor() Processor {}
// Reset重置清除整个区块链，将其恢复到genesis state.
func (bc *BlockChain) Reset() error {}
 
// ResetWithGenesisBlock 清除整个区块链, 用特定的genesis state重塑，被Reset所引用
func (bc *BlockChain) ResetWithGenesisBlock(genesis *types.Block) error {}
 
// repair尝试通过回滚当前块来修复当前的区块链，直到找到具有关联状态的块。
// 用于修复由崩溃/断电或简单的非提交尝试导致的不完整的数据库写入。
//此方法仅回滚当前块。 当前标头和当前快速块保持不变。
func (bc *BlockChain) repair(head **types.Block) error {}
 
// reorgs需要两个块、一个旧链以及一个新链，并将重新构建块并将它们插入到新的规范链中，并累积潜在的缺失事务并发布有关它们的事件
func (bc *BlockChain) reorg(oldBlock, newBlock *types.Block) error{}
 
// Rollback 旨在从数据库中删除不确定有效的链片段
func (bc *BlockChain) Rollback(chain []common.Hash) {}
// SetReceiptsData 计算收据的所有非共识字段
func SetReceiptsData(config *params.ChainConfig, block *types.Block, receipts types.Receipts) error {}
 
// SetHead将本地链回滚到指定的头部。
// 通常可用于处理分叉时重选主链。对于Header，新Header上方的所有内容都将被删除，新的头部将被设置。
// 但如果块体丢失，则会进一步回退（快速同步后的非归档节点）。
func (bc *BlockChain) SetHead(head uint64) error {}
 
// SetProcessor设置状态修改所需要的processor
func (bc *BlockChain) SetProcessor(processor Processor) {}
 
// SetValidator 设置用于验证未来区块的validator
func (bc *BlockChain) SetValidator(validator Validator) {}
 
// State 根据当前头区块返回一个可修改的状态
func (bc *BlockChain) State() (*state.StateDB, error) {}
 
// StateAt 根据特定时间点返回新的可变状态
func (bc *BlockChain) StateAt(root common.Hash) (*state.StateDB, error) {}
 
// Stop 停止区块链服务，如果有正在import的进程，它会使用procInterrupt来取消。
// it will abort them using the procInterrupt.
func (bc *BlockChain) Stop() {}
 
// TrieNode从memory缓存或storage中检索与trie节点hash相关联的数据。
func (bc *BlockChain) TrieNode(hash common.Hash) ([]byte, error) {}
 
// Validator返回当前validator.
func (bc *BlockChain) Validator() Validator {}
 
// WriteBlockWithoutState仅将块及其元数据写入数据库，但不写入任何状态。 这用于构建竞争方叉，直到超过规范总难度。
func (bc *BlockChain) WriteBlockWithoutState(block *types.Block, td *big.Int) (err error){}
 
// WriteBlockWithState将块和所有关联状态写入数据库。
func (bc *BlockChain) WriteBlockWithState(block *types.Block, receipts []*types.Receipt, state *state.StateDB) {}
 
// writeHeader将标头写入本地链，因为它的父节点已知。 如果新插入的报头的总难度变得大于当前已知的TD，则重新路由规范链
func (bc *BlockChain) writeHeader(header *types.Header) error{}
 
// 处理未来区块链
func (bc *BlockChain) update() {}
```

## blockchain初始化(NewBlockChain)

主要步骤：

①：创建一个新的headerChain结构

```go
bc.hc, err = NewHeaderChain(db, chainConfig, engine, bc.getProcInterrupt)
```

1. 根据**number（0）**获取**genesisHeader**
2. 从**rawdb中读取HeadBlock并存储在currentHeade**r中

②：获取genesisBlock

```go
bc.genesisBlock = bc.GetBlockByNumber(0)
```

③：如果链不为空，则用老的链数据初始化链

```go
if bc.empty() {
		rawdb.InitDatabaseFromFreezer(bc.db)
	}
```

④：加载最新的状态数据

```go
if err := bc.loadLastState(); err != nil {
		return nil, err
	}
```

⑤：检查区块哈希的当前状态，并确保链中没有任何坏块

```go
for hash := range BadHashes {
		if header := bc.GetHeaderByHash(hash); header != nil {
			headerByNumber := bc.GetHeaderByNumber(header.Number.Uint64())
			if headerByNumber != nil && headerByNumber.Hash() == header.Hash() {
				log.Error("Found bad hash, rewinding chain", "number", header.Number, "hash", header.ParentHash)
				bc.SetHead(header.Number.Uint64() - 1)
				log.Error("Chain rewind was successful, resuming normal operation")
			}
		}
	}
```

⑥：定时处理future block

```go
go bc.update()
	->procFutureBlocks
		->InsertChain
```

总的来说做了以下几件事：

1. 配置cacheConfig，创建各种lru缓存
2. 初始化triegc
3. 初始化stateDb：state.NewDatabase(db)
4. 初始化区块和状态验证：NewBlockValidator()
5. 初始化状态处理器：NewStateProcessor()
6. 初始化区块头部链：NewHeaderChain()
7. 查找创世区块：bc.genesisBlock = bc.GetBlockByNumber(0)
8. 加载最新的状态数据：bc.loadLastState()
9. 检查区块哈希的当前状态，并确保链中没有任何坏块
10. go bc.**futureBlocksLoop**()定时处理future block

## 加载区块链状态(locadLastState)

1：从rawdb数据库中恢复**最新的**headblock，如果rawdb数据库空的话也就是没有读出来头部区块的话，触发reset chain

```go
//获取最新的头部hash
head := rawdb.ReadHeadBlockHash(bc.db)
	if head == (common.Hash{}) {
		log.Warn("Empty database, resetting chain")
		return bc.Reset()
	}
```

2：通过头部hash获取头部区块，确保整个head block是可以获取的，若为空，则触发reset chain

```go
//通过头部hash获取头部区块
currentBlock := bc.GetBlockByHash(head)
	if currentBlock == nil {
		// Corrupt or empty database, init from scratch
		log.Warn("Head block missing, resetting chain", "hash", head)
		return bc.Reset()
	}
```

其中GetBlockByHash(head)的方法是这样的：

```go
func (bc *BlockChain) GetBlockByHash(hash common.Hash) *types.Block {
    //通过hash获取number
	number := bc.hc.GetBlockNumber(hash)
	if number == nil {
		return nil
	}
    //通过number获取区块
	return bc.GetBlock(hash, *number)
}
```

3：存储当前的headblock和设置当前的headHeader以及头部快速块

```go
bc.currentBlock.Store(currentBlock)
....
bc.hc.SetCurrentHeader(currentHeader)
...
bc.currentFastBlock.Store(currentBlock)
```

具体的`store`函数的使用我写了一段代码作为示例帮助理解：

```go
package main

import (
	"fmt"
	"sync/atomic"
)

func main() {
	var res atomic.Value
	res.Store("test1")
	fmt.Println(res)//res=>{test1}

	res.Store("test2")
	fmt.Println(res)//res=>{test}
}
```

证明里面存放的数据会进行覆盖，所以类似于这样的代码就好理解了：

```go
bc.currentFastBlock.Store(currentBlock)
	headFastBlockGauge.Update(int64(currentBlock.NumberU64()))

	//ReadHeadFastBlockHash retrieves the hash of the current fast-sync head block.
	if head := rawdb.ReadHeadFastBlockHash(bc.db); head != (common.Hash{}) {
		if block := bc.GetBlockByHash(head); block != nil {
			//为何store完之后又进行一次store 
            //其实是为了和链进行同步（其他地方链可能
            bc.currentFastBlock.Store(block)
			headFastBlockGauge.Update(int64(block.NumberU64()))
		}
	}
```



## 插入数据到blockchain中

①：如果链正在中断，直接返回

②：开启并行的签名恢复

③：校验header

```go
abort, results := bc.engine.VerifyHeaders(bc, headers, seals)
```

④：循环校验body

```go
block, err := it.next()
	-> ValidateBody
		-> VerifyUncles
```

包括以下错误：

- block已知
- uncle太多
- 重复的uncle
- uncle是祖先块
- uncle哈希不匹配
- 交易哈希不匹配
- 未知祖先
- 祖先块的状态无法获取

如果block存在，且是已知块，则写入已知块。

如果是祖先块的状态无法获取的错误，则作为侧链插入：

```go
bc.insertSideChain(block, it)
```

如果是未来块或者未知祖先，则添加未来块：

```go
bc.addFutureBlock(block);
```

如果是其他错误，直接中断，并且报告坏块。

```go
bc.futureBlocks.Remove(block.Hash())
...
bc.reportBlock(block, nil, err)
```

⑤：没有校验错误

如果是坏块，则报告；如果是未知块，则写入未知块；根据给定trie，创建state；

执行块中的交易：

```go
receipts, logs, usedGas, err := bc.processor.Process(block, statedb, bc.vmConfig)
```

使用默认的validator校验状态：

```go
bc.validator.ValidateState(block, statedb, receipts, usedGas);
```

将块写入到区块链中并获取状态：

```go
status, err := bc.writeBlockWithState(block, receipts, logs, statedb, false)
```

⑥：校验写入区块的状态

- CanonStatTy ： 插入成功新的block
- SideStatTy：插入成功新的分叉区块
- Default：插入未知状态的block

⑦：如果还有块，并且是未来块的话，那么将块添加到未来块的缓存中去

```go
bc.addFutureBlock(block)
```

至此insertChain 大概介绍清楚。

------

## 将块和关联状态写入到数据库

函数：**WriteBlockWithState**

①：计算父块的total td

```go
ptd := bc.GetTd(block.ParentHash(), block.NumberU64()-1)
```

②：添加待插入块本身的td ,并将此时最新的total td 存储到数据库中。

```GO
bc.hc.WriteTd(block.Hash(), block.NumberU64(), externTd)
```

③：将块的header和body分别序列化到数据库

```go
rawdb.WriteBlock(bc.db, block)
	->WriteBody(db, block.Hash(), block.NumberU64(), block.Body())
	->WriteHeader(db, block.Header())
```

④：将状态写入底层内存Trie数据库

```go
state.Commit(bc.chainConfig.IsEIP158(block.Number()))
```

⑤：存储一个块的所有交易数据

```go
rawdb.WriteReceipts(batch, block.Hash(), block.NumberU64(), receipts)
```

⑥：将新的head块注入到当前链中

```go
if status == CanonStatTy {
		bc.insert(block)
	}
```

- 存储分配给规范块的哈希
- 存储头块的哈希
- 存储最新的快
- 更新currentFastBlock

到此writeBlockWithState 结束，从上面可以知道，insertChain的最终还是调用了writeBlockWithState的insert方法完成了最终的插入动作。

------

## 思考

1. 为什么还要导入已知块？？？writeKnownBlock

参考：

> https://github.com/mindcarver/blockchain_guide (优秀的区块链学习营地)

