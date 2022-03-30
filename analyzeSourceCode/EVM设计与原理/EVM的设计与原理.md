# 前言

在阅读这篇文章之前，请您先阅读[初步理解以太坊虚拟机](https://github.com/learnerLj/geth-analyze/blob/main/analyzeSourceCode/EVM%E8%AE%BE%E8%AE%A1%E4%B8%8E%E5%8E%9F%E7%90%86/%E5%88%9D%E6%AD%A5%E7%90%86%E8%A7%A3%E4%BB%A5%E5%A4%AA%E5%9D%8A%E8%99%9A%E6%8B%9F%E6%9C%BA.md)，它将会介绍 EVM 的基本知识。在开始之前，假设您已经掌握了上文中的基础，我们根据黄皮书进一步地补充理论基础。由于原始的黄皮书公式过多，不易阅读，可以参考按照论文[重写后的版本](https://github.com/chronaeon/beigepaper)。其次，本文使用的图片来自其他资料，会在参考资料部分注明。

# EVM 设计原理

以太坊可以抽象的分成两部分，一部分是状态，另外一部分是用于改变状态的 EVM。因此，以太坊在整体上可以看作一个**基于交易的状态机**：起始于一个创世区块（Genesis）状态，然后随着交易的执行，状态逐步改变。

关于**不可篡改性**，黄皮书的重写版的表述不错：

> Ethereum programs can be trusted to execute without any interference from external non-network forces.

## 基本概念

**以太币**：以太坊发行自身的货币，用于衡量计算消耗，它不仅是作为金融工具，更是要作为世界计算机，为所有的应用服务。所有的交易在机器层面都是以 wei 作为单位。

**世界状态**：世界状态可以被视作以太坊地址到账户状态的映射。存储时，地址和值经过 RLP 编码，以键值对的形式，通过 MPT 的组织方式存储在数据库中。这个数据库被称作**状态数据库**。

**MPT**：前辈们的分析飞铲精湛，建议仔细阅读我们整合、修正过的 [MPT树](https://github.com/learnerLj/geth-analyze/blob/main/analyzeSourceCode/MPT%E6%A0%91.md).

RLP：



## 交易的执行



### 运行前测试

交易的执行是以太坊协议中最复杂的部分。首先**任意交易在执行之前必须通过初始的有效性测试**。包括：

1. 交易是 RLP 格式数据，没有多余的后缀字节； 

2. 交易的签名是有效的；
3.  交易的 nonce 是有效的（等于发送者账户当前的 nonce）；
4.  gas 上限不小于交易所要使用的 gas;
5. 发送者账户的 balance 应该不少于实际费用，且需要提前支付。

### 子状态

交易的执行过程中会累积产生一些特定的信 息，我们称为交易子状态，它包括四部分：

1. 自毁集合：一组应该在交易完成后被删除的账户。
2. 一系列的日志：这是一些列针归档的、可索引的“检查点”，允许以太坊世界的外部旁观者（例如去中心化应用的前端）来简单地跟踪合约调用。
3. 交易所接触过的账户集合，其中的空账户可以在交易结束时删除。
4. 最后是应该返还的余额；

### 合约创建

创建合约需要的参数将会在源码分析中给出。心的





## 架构解析

![以太坊虚拟机 (EVM) 架构和执行上下文](https://cypherpunks-core.github.io/ethereumbook/images/evm-architecture.png)

## 











































# 参考

- 以太坊黄皮书
- https://cypherpunks-core.github.io/ethereumbook/13evm.html
- https://github.com/chronaeon/beigepaper

