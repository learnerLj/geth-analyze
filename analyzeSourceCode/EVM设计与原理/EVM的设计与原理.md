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

### **以太币**：

以太坊发行自身的货币，用于衡量计算消耗，它不仅是作为金融工具，更是要作为世界计算机，为所有的应用服务。所有的交易在机器层面都是以 wei 作为单位。

### **世界状态**：

世界状态可以被视作以太坊地址到账户状态的映射。存储时，地址和值经过 RLP 编码，以键值对的形式，通过 MPT 的组织方式存储在数据库中。这个数据库被称作**状态数据库**。

![image-20220330102850761](http://blog-blockchain.xyz/202203301028818.png)

### **MPT**：

前辈们的分析非常精湛，建议仔细阅读我们整合、修正过的 [MPT树](https://github.com/learnerLj/geth-analyze/blob/main/analyzeSourceCode/MPT%E6%A0%91.md).

### **RLP**：

我们在[以太坊的数据组织](https://github.com/learnerLj/geth-analyze/blob/main/analyzeSourceCode/EVM%E8%AE%BE%E8%AE%A1%E4%B8%8E%E5%8E%9F%E7%90%86/%E4%BB%A5%E5%A4%AA%E5%9D%8A%E7%9A%84%E6%95%B0%E6%8D%AE%E7%BB%84%E7%BB%87.md)中介绍了 RLP 的编码规则，适合对数据组织的方式形成基本的认识。

### **账户状态**：

前面提到了状态是以键值对的方式存储，账户状态也是以键值对存储，主要包括：nonce、balance、storage root(256 位的的账户数的 MPT 树根)、code hash(字节码的哈希，当账户接收到消息后改变)

![image-20220330103327705](http://blog-blockchain.xyz/202203301033775.png)

### **布隆过滤器**：

我们尝试解读过，但是由于缺乏工程经验，对于并发和调度不熟悉，因此只是半成品，但是也有一定的参考意义。可见[这篇文章](https://github.com/learnerLj/geth-analyze/blob/main/analyzeSourceCode/%E5%B8%83%E9%9A%86%E8%BF%87%E6%BB%A4%E5%99%A8.md)。

### **区块**：

请阅读这篇[笔记](https://github.com/learnerLj/geth-analyze/blob/main/analyzeSourceCode/%E5%9F%BA%E6%9C%AC%E6%95%B0%E6%8D%AE%E7%BB%93%E6%9E%84/%E7%90%86%E8%A7%A3%E5%8C%BA%E5%9D%97.md)。

### **收据**：

收据可以用于索引、零知识证明等方面，它是交易执行中某些信息的编码。详细内容请阅读博客——[理解收据](https://github.com/learnerLj/geth-analyze/blob/main/analyzeSourceCode/%E5%9F%BA%E6%9C%AC%E6%95%B0%E6%8D%AE%E7%BB%93%E6%9E%84/%E7%90%86%E8%A7%A3%E6%94%B6%E6%8D%AE.md)。

## 交易的处理和执行

交易是以太坊账户之间通信的最基本的方式，可以视作是签名后发送给 EVM 的执行指令。每一笔交易都会造成以太坊状态的改变以及产生临时存储的状态。

![image-20220402192210613](http://blog-blockchain.xyz/202204021922731.png)

交易的组成如下

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

交易费的收取可以分成三部分，**第一部分**是最普遍的基础费用和指令消耗。**第二部分**是用于交易中的子消息调用或者创建合约。**第三类**是内存拓展的开销。当交易执行时，需要多少内存并不是预定的，而是根据操作的需要，拓展 32 字节一组的 slot。存储是抱着能省就省的目的，因此清除存储中的某一项内容，不但不消耗交易费，反而会退回一部分手续费。

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

消息调用很类似于交易的执行，黄皮书上定义**合约账户接到指令后调用其他对象的行为是消息调用**。执行时会创建运行时的子对象，每个子对象对应子状态，必须等所有的子对象都完成计算后，才能确定最终的状态。

消息调用的环境参数包括：

- account_address：当前上下文字节码所在的账户。
- sender_address：发起消息调用的账户。
- originator_price：最初确定的 gasPrice
- input_data：确定调用行为的指令数据，例如执行那个函数，输入什么参数。
- newstate_value：消息调用的转账金额。
- code：将要执行的机器码，`account_address` 所有。
- block_header：当前消息调用的区块头。
- stack_depth：从最开始的执行计算，执行路径上 call 指令的底层调用的次数，也就是调用的堆栈深度。

![../_images/Picture50.png](https://ethbook.abyteahead.com/_images/Picture50.png)

注意创建合约的交易和调用合约的交易，在处理上是非常不一样的。创建合约的交易的 `to` 字段为空。

![image-20220402192357434](http://blog-blockchain.xyz/202204021923501.png)

总而言之，交易的执行可以抽象成运行的执行状态和系统状态逐步改变的过程。运行时的执行状态叫做 `machine_state` ，它包括：

- 可使用的 gas
- 程序计数器 PC
- 内存的内容
- 使用的内存的长度
- 栈的内容

具体的执行过程，请阅读 [智能合约审计的深入字节码分析部分](https://github.com/learnerLj/geth-analyze/blob/main/analyzeSourceCode/%E6%99%BA%E8%83%BD%E5%90%88%E7%BA%A6%E5%88%86%E6%9E%90%E5%9F%BA%E7%A1%80/%E6%99%BA%E8%83%BD%E5%90%88%E7%BA%A6%E5%AE%A1%E8%AE%A1.md#%E6%B7%B1%E5%85%A5%E5%AD%97%E8%8A%82%E7%A0%81%E5%88%86%E6%9E%90)。

在消息调用时有几个预编译的合约，作为链的基础架构的一部分。

详细可见 https://www.evm.codes/precompiled

![image-20220402104011700](http://blog-blockchain.xyz/202204021040002.png)

### 执行终止

**交易具有原子性**，要么完全的执行成功，如果执行过程某部分执行异常，会立即注销状态。异常终止并不是通过操作码实现，而是通过一系列的检查完成的，有以下情况：

- gas 不足。 
- 指令错误
- 栈溢出，超过 1024 的限制
- PC 跳转错误。
- 静态调用中尝试修改状态。

![image-20220402193054481](http://blog-blockchain.xyz/202204021930555.png)

具体的说，

1. 每执行一条指令，就扣除相应的 gas。这里从 `machine_status` 中获取当前 gas。
2. 对每个循环附带布尔标志。true 表示循环时异常终止；false 表示正常完成，如果直到需要执行的交易的集合全部转换成栈操作的序列时，仍然保持 false，表示可以通过控制流正常终止。

正常终止可以分成两种情况：

1. RETURN 或者 REVERT，二者都会在终止时会执行特殊的终止函数。
2. STOP 或者 SELFDESTRUCT，他们会直接终止，销毁执行时的状态。

## 区块的形成

以太坊的区块链可以视作是一棵从树根到叶子的区块树，分叉是树的分叉，主网选择区块树的工作量最大的路径以保持共识。一般而言，这条路径上的叶子是最多的，每个叶子对应一个成功验证的有效区块。路径越长，挖矿所需要的努力就越多。

形成区块的流程如下：

1. 验证叔块。每个叔块的区块头必须是有效的，而且需要在前 6 个区块以内。至多引用两个叔块。
2. 验证交易。区块的 gasUsed 必须和交易的 gasUsed 之和相等。
3. 发放挖矿奖励。如果引用叔块，那么叔块的挖矿者和当前块的挖矿者都会得到奖励。如果被引用的叔块 A 引用了叔块 B，而 A、B 都被当前区块引用，那么 B 的奖励可以叠加。
4. 校验状态和区块 nonce.









## 架构解析

![以太坊虚拟机 (EVM) 架构和执行上下文](https://cypherpunks-core.github.io/ethereumbook/images/evm-architecture.png)

![img](https://pic2.zhimg.com/80/v2-90964a7f9d10855d21287c768875068d_1440w.jpg)











































# 参考

- https://ethereum.github.io/yellowpaper/paper.pdf ，也有[中文翻译版本](https://github.com/riversyang/ethereum_yellowpaper)，但是版本较旧。黄皮书在提出后经过许多次的修订。
- https://cypherpunks-core.github.io/ethereumbook/13evm.html
- 黄皮书的[论文重写版本](https://github.com/chronaeon/beigepaper)，但是它的语言写的不易懂，为了信息量写了很多很长的句子，句子结构也很复杂，读起来可能还不如原版舒服。
- https://takenobu-hs.github.io/downloads/ethereum_evm_illustrated.pdf
- https://ethbook.abyteahead.com/ch7/call.html

