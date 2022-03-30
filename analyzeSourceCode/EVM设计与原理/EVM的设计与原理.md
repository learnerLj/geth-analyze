# 前言

在阅读这篇文章之前，请您先阅读[初步理解以太坊虚拟机](https://github.com/learnerLj/geth-analyze/blob/main/analyzeSourceCode/EVM%E8%AE%BE%E8%AE%A1%E4%B8%8E%E5%8E%9F%E7%90%86/%E5%88%9D%E6%AD%A5%E7%90%86%E8%A7%A3%E4%BB%A5%E5%A4%AA%E5%9D%8A%E8%99%9A%E6%8B%9F%E6%9C%BA.md)和[以太坊的数据组织](https://github.com/learnerLj/geth-analyze/blob/main/analyzeSourceCode/EVM%E8%AE%BE%E8%AE%A1%E4%B8%8E%E5%8E%9F%E7%90%86/%E4%BB%A5%E5%A4%AA%E5%9D%8A%E7%9A%84%E6%95%B0%E6%8D%AE%E7%BB%84%E7%BB%87.md)，它将会介绍 EVM 的基本知识，帮助您形成基本的认识。在开始之前，假设您已经掌握了上文中的基础，我们根据黄皮书进一步地补充理论基础。由于原始的黄皮书公式过多，不易阅读，可以参考按照论文[重写后的版本](https://github.com/chronaeon/beigepaper)。其次，本文使用的图片来自其他资料，会在参考资料部分注明。

非常推荐读者观看这个视频：[EVM: From Solidity to byte code, memory and storage](https://www.youtube.com/watch?v=RxL_1AfV7N4&t=1s) ，这是配套的 [PDF](https://slack-files.com/T9C7VSRBN-F0154NTUM3L-3eefe73def)。

它梳理源码到字节码的流程，演示操作码的变化，非常棒。如果读者通过前面提到的文章以及理解 EVM 的存储空间布局的话，这个视频可以为您提供字节码编写合约的基础:smirk_cat:.请善用 Remix IDE 的单步调试功能，可以通过实操，大大加深理解。

# EVM 设计原理

以太坊可以抽象的分成两部分，一部分是状态，另外一部分是用于改变状态的 EVM。因此，以太坊在整体上可以看作一个**基于交易的状态机**：起始于一个创世区块（Genesis）状态，然后随着交易的执行，状态逐步改变。

![image-20220330102802873](http://blog-blockchain.xyz/202203301028001.png)

关于**不可篡改性**，黄皮书的重写版的表述不错：

> Ethereum programs can be trusted to execute without any interference from external non-network forces.

> Rather than storing program code in generally-accessible memory or storage, it is stored separately in a virtual ROM interactable only through specialized instructions

不允许外部的任何干扰，特殊的读取程序的方式。

## 基本概念

**以太币**：以太坊发行自身的货币，用于衡量计算消耗，它不仅是作为金融工具，更是要作为世界计算机，为所有的应用服务。所有的交易在机器层面都是以 wei 作为单位。

**世界状态**：世界状态可以被视作以太坊地址到账户状态的映射。存储时，地址和值经过 RLP 编码，以键值对的形式，通过 MPT 的组织方式存储在数据库中。这个数据库被称作**状态数据库**。

![image-20220330102850761](http://blog-blockchain.xyz/202203301028818.png)

**MPT**：前辈们的分析飞铲精湛，建议仔细阅读我们整合、修正过的 [MPT树](https://github.com/learnerLj/geth-analyze/blob/main/analyzeSourceCode/MPT%E6%A0%91.md).

**RLP**：我们在[以太坊的数据组织](https://github.com/learnerLj/geth-analyze/blob/main/analyzeSourceCode/EVM%E8%AE%BE%E8%AE%A1%E4%B8%8E%E5%8E%9F%E7%90%86/%E4%BB%A5%E5%A4%AA%E5%9D%8A%E7%9A%84%E6%95%B0%E6%8D%AE%E7%BB%84%E7%BB%87.md)中介绍了 RLP 的编码规则，适合对数据组织的方式形成基本的认识。

**账户状态**：前面提到了状态是以键值对的方式存储，账户状态也是以键值对存储，主要包括：nonce、balance、storage root(256 位的的账户数的 MPT 树根)、code hash(字节码的哈希，当账户接收到消息后改变)

![image-20220330103327705](http://blog-blockchain.xyz/202203301033775.png)

**布隆过滤器**：我们尝试解读过，但是由于缺乏工程经验，对于并发和调度不熟悉，因此只是半成品，但是也有一定的参考意义。可见[这篇文章](https://github.com/learnerLj/geth-analyze/blob/main/analyzeSourceCode/%E5%B8%83%E9%9A%86%E8%BF%87%E6%BB%A4%E5%99%A8.md)。

**区块**：请阅读这篇[笔记](https://github.com/learnerLj/geth-analyze/blob/main/analyzeSourceCode/%E5%9F%BA%E6%9C%AC%E6%95%B0%E6%8D%AE%E7%BB%93%E6%9E%84/%E7%90%86%E8%A7%A3%E5%8C%BA%E5%9D%97.md)。

**收据**：收据可以用于索引、零知识证明等方面，它是交易执行中某些信息的编码。详细内容请阅读博客——[理解收据](https://github.com/learnerLj/geth-analyze/blob/main/analyzeSourceCode/%E5%9F%BA%E6%9C%AC%E6%95%B0%E6%8D%AE%E7%BB%93%E6%9E%84/%E7%90%86%E8%A7%A3%E6%94%B6%E6%8D%AE.md)。

## 交易的处理和执行

交易是以太坊账户之间通信的最基本的方式，可以视作是签名后发送给 EVM 的执行指令。每一笔交易都会造成以太坊状态的改变以及产生临时存储的状态。交易的组成如下

- nonce: 由交易发送者发出的的交易的数量。
- gasPrice: 每单位 gas 的价格
- gasLimit: 用于执行这个交易的最大 gas 数量。
- to: 160 位的消息调用接收者地址；创建合约时为 0 地址。
- value: 转移到接收者账户的 Wei 的数量；如果是创建合约，则代表给新建合约地址的初始余额。
- v, r, s: 与交易签名相符的若干数值，用于确定交易 的发送者。
- init： 一个不限制大小的字节数组，用来指定账户初始化程序的 EVM 代码。它仅会在合约创建时被执行一 次，然后就会被丢弃，不会存在链上。
- data: 一个不限制大小的字节数组，用来指定消息调用的输入数据

![image-20220330165714893](http://blog-blockchain.xyz/202203301657199.png)

请注意，这是最开始的设计思路，后面经过诸多的 EIP 后，有些改变，更详细的内容可见 [理解交易](https://github.com/learnerLj/geth-analyze/blob/main/analyzeSourceCode/%E5%9F%BA%E6%9C%AC%E6%95%B0%E6%8D%AE%E7%BB%93%E6%9E%84/%E7%90%86%E8%A7%A3%E4%BA%A4%E6%98%93.md)。

### 交易费的设置

交易费的收取可以分成三部分，**第一部分**是最普遍的执行计算消耗。**第二部分**是用于交易中的子消息调用或者创建合约。**第三类**是内存拓展的开销。当交易执行时，需要多少内存并不是预定的，而是根据操作的需要，拓展 32 字节一组的 slot。存储是抱着能省就省的目的，因此清除存储中的某一项内容，不但不消耗交易费，反而会退回一部分手续费。

### 运行前检查

交易的执行是以太坊协议中最复杂的部分。首先**任意交易在执行之前必须通过初始的有效性测试**。包括：

1. 交易是 RLP 格式数据，没有多余的后缀字节； 

2. 交易的签名是有效的；
3.  交易的 nonce 是有效的（等于发送者账户当前的 nonce）；
4.  gas 上限不小于交易所要使用的 gas;
5. 发送者账户的 balance 应该不少于实际费用，且需要提前支付。

### 子状态和临时状态

交易的执行过程中会累积产生一些特定的信 息，我们称为交易子状态，它包括四部分：

1. 自毁集合：一组应该在交易完成后被删除的账户。
2. 一系列的日志：这是一些归档的、可索引的“检查点”，允许在区块链以外，简单地跟踪合约调用。
3. 交易所接触过的账户集合，其中的空账户可以在交易结束时删除。
4. 最后是应该返还的余额；

在交易执行后，确定最终的状态前，会处于一个临时状态，在这个状态中，操作码会逐步执行。这个状态包括：

- 临时状态
- 操作码计算开销
- 相关子状态
- 在交易完成后，执行 selfdestruct 时，接受合约余额账户。
- 一系列日志，包括布隆过滤器和收据，用于外部应用的跟踪执行过程和查询。
- 将返回的余额（例如gas剩余、销毁操作返回。

如果交易失败，这些状态将会重置为空，因此达到了回滚的目的，执行失败的交易不影响世界状态。因此，为了改变世界状态，交易要么完整的执行完毕，要么毫无作用。

### 交易的收据

前面提到的子状态中包括日志，日志中有比较特殊的一项，叫做交易的收据，**用于记录交易的执行结果**。日志的集合与包含事件的布隆过滤器，都存储在收据中。交易执行后的状态码和使用的 gas 也在收据中。

为了方便读者理解，下面是 geth 中的定义：

```go
type Receipt struct {
   // Consensus fields: These fields are defined by the Yellow Paper
   Type      uint8  `json:"type,omitempty"` //交易类型
   PostState []byte `json:"root"`           //交易成功/失败时的 RLP 编码
   Status    uint64 `json:"status"`         //交易成功/失败的状态码
   //区块中直到这一笔交易累积使用的 gas
   CumulativeGasUsed uint64 `json:"cumulativeGasUsed" gencodec:"required"`
   //布隆过滤器
   Bloom Bloom `json:"logsBloom"         gencodec:"required"`
   //合约的日志列表
   Logs []*Log `json:"logs"              gencodec:"required"`

   //处理交易时的字段

   // Implementation fields: These fields are added by geth when processing a transaction.
   // They are stored in the chain database.
   TxHash          common.Hash    `json:"transactionHash" gencodec:"required"`
   ContractAddress common.Address `json:"contractAddress"`
   GasUsed         uint64         `json:"gasUsed" gencodec:"required"`

   //记录区块信息和交易索引，用于叫检查交易与对应收据的兼容性
   // Inclusion information: These fields provide information about the inclusion of the
   // transaction corresponding to this receipt.
   BlockHash        common.Hash `json:"blockHash,omitempty"`
   BlockNumber      *big.Int    `json:"blockNumber,omitempty"`
   TransactionIndex uint        `json:"transactionIndex"`
}
```

### 消息调用









### 合约创建

创建合约需要的参数将会在源码分析中给出。心的





## 架构解析

![以太坊虚拟机 (EVM) 架构和执行上下文](https://cypherpunks-core.github.io/ethereumbook/images/evm-architecture.png)













































# 参考

- 以太坊黄皮书
- https://cypherpunks-core.github.io/ethereumbook/13evm.html
- https://github.com/chronaeon/beigepaper
- https://takenobu-hs.github.io/downloads/ethereum_evm_illustrated.pdf

