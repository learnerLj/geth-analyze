# 理论基础

在阅读这篇文章之前，请您先阅读[初步理解以太坊虚拟机](https://github.com/learnerLj/geth-analyze/blob/main/analyzeSourceCode/EVM%E8%AE%BE%E8%AE%A1%E4%B8%8E%E5%8E%9F%E7%90%86/%E5%88%9D%E6%AD%A5%E7%90%86%E8%A7%A3%E4%BB%A5%E5%A4%AA%E5%9D%8A%E8%99%9A%E6%8B%9F%E6%9C%BA.md)，它将会介绍 EVM 的基本知识。在开始之前，假设您已经掌握了上文中的基础，我们根据黄皮书进一步地补充理论基础。由于原始的黄皮书公式过多，不易阅读，可以参考按照论文[重写后的版本](https://github.com/chronaeon/beigepaper)。

## 交易的执行

以太坊在整体上可以看作一个**基于交易的状态机**：起始 于一个创世区块（Genesis）状态，然后随着交易的执行状 态逐步改变一直到最终状态，这个最终状态是以太坊世界 的权威“版本”。

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

