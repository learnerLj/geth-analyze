> 作者：[知乎-山上的石头](https://zhuanlan.zhihu.com/p/459969916)

> GitHub: **[geth-analyze](https://github.com/learnerLj/geth-analyze)**

# 前言

这是我学习区块链入门时做的笔记（基于 0.8.7 版本），基本涵盖了编写合约所需常用知识，由于做智能合约安全方面的研究需要精通 Solidity 和 以太坊原理，因此做的笔记比较详实。

这些笔记基于阅读英文文档，参考中文文档和 stack overflow 以及相关教程，我根据学习者的接受新知识的顺序，对文章结构做了适当优化。

这篇文章既可以作为新手入门（因为啃英文文档并且搜索信息并不是容易的事情），也可以作为快捷的检索帮助文档（官方文档的翻译某些部分比较难以理解）。建议使用电脑端阅读。

初稿完成时，都还没学 C 语言，只是一知半解的边学边记。在大二上学期的寒假，我重新整理了一遍，修正了部分错误，将拗口的表述转化成习惯表述，补充了文档中缺少的范例，根据经验突出需要强调的注意事项，使得读者可以跟容易的学习。

本文共计接近7万字，如果觉得有帮助点赞关注呀，我将会继续写智能合约的攻击方式、以太坊虚拟机原理、字节码的深入探索等等，逐渐完善知识体系，并且会分享读论文时的前沿理论。

参考：

- [Solidity 最新(0.8.0)中文文档](https://link.zhihu.com/?target=https%3A//learnblockchain.cn/docs/solidity/index.html)
- [Solidity - Solidity 0.8.12 documentation](https://link.zhihu.com/?target=https%3A//docs.soliditylang.org/en/latest/index.html)
- [https://solidity-by-example.org](https://link.zhihu.com/?target=https%3A//solidity-by-example.org/)



# 智能合约介绍

## 智能合约架构

合约自底向上可以分为基础设施层、合约层、运维层、智能层、表现层、应用层。如图所示

![image-20220111202432714](http://blog-blockchain.xyz/202203260119810.png)

- 基础设施层主要是合约的可信执行环境，包括共识机制、激励机制等等。

- 合约层封装了静态的智能合约对象，包括合约对象、只读调用方式对象、只写调用方式对象。详细一点地说，有满足条件时触发地响应规则，对外暴露的接口等。

- 运维层封装了对静态的合约数据的动态操作，例如形式化验证、安全审查、销毁等。

- 智能层主要封装了各类为了满足业务需求的各类算法。虽然现在的智能合约具有局限性，计算密集型任务难以胜任，但是随着认知计算、小样本学习、合约模型优化等技术的发展，将会出现一定的智能性。

- 表现层主要是合约在具体应用中的表现形式，合约类似于应用的接口，它可以用于各类业务场景。

- 应用层主要是合约在具体领域的应用。



**运行机制大致如下**

![image-20220111204520558](http://blog-blockchain.xyz/202203260119275.png)

EVM 是以太坊虚拟机，也是状态机，通过 oracles 预言机获取链外信息，然后智能合约编写”如果...就执行...“的算法，改变 EVM 的状态。

链外的用户也可以通过交易调用合约，它们也可以作为矿工生成区块，区块生成后通过共识机制广播到其他节点，最终确认后形成区块。

## 合约的执行环境

在以太坊中，EVM 是合约的执行环境。操作系统上安装 geth 之类的区块链节点客户端，客户端维持 EVM, EVM 维护着可信的执行环境。

EVM 是无寄存器、基于栈的虚拟机，他的存储空间有三类，stack、memory、storage。stack 和 memory 都是临时存储，在智能合约运行时有效，当运行结束后回收；storage 是永久性存储，因此它的操作更消耗计算资源，gas 也高得多。

stack 是运行时必须的资源，以32字节为一组访问（以太坊是 256 位的虚拟机），Solidity 规定不能够超过 1024层。memory 主要是临时存储数组、字符串等较大的数据，以1个字节为1组，更加灵活。这样的储存空间分类，同样在 Solidity 中有所体现。

智能合约的静态数据存储在虚拟机的状态中，包含函数选择器和函数入口，当调用合约代码的某个函数时，就会通过函数签名匹配入口和参数。

## 合约的属性

### gas机制

gas 用于衡量每一项操作消耗的计算资源，如果 gaslimit 小于消耗的 gas 的话，执行会失败，状态回滚到调用之前，消耗的 gas 不会退还。

### 异常传递机制

对合约的调用可以分为内部调用和外部调用。内部调用只需要在 EVM 中跳转到对应位置即可，效率更高。而外部调用需要用 CALL 向其他合约发送消息，消息的格式需要满足外部合约的接口规范。**外部调用时执行失败，不会传递到当前合约，指挥返回布尔值**。这也是造成重入攻击的原因。

### 委托调用

DELEGATECALL 指令也可以调用外部合约的函数，它的特点是调用外部合约的函数时，这个函数的上下文信息是当前合约的，也即当前合约将自身状态注入被调用的合约。这也造成和风险，因为一旦拥有当前合约的上下文状态，那么就很容易做坏事。

### 合约无法修改

代码一经转换成 ABI 和字节码，部署上链后，就无法修改。

### 调用序列

每个合约有自己的全局变量和状态，一个交易可能触发一整条调用序列。例如A合约调用B合约和C合约，B合约调用C合约，C合约调用D合约和A合约，这样复杂的调用会改变参与的合约的状态。这样的执行顺序给静态分析带来了许多挑战。

















