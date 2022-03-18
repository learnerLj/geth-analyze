---
title: txnoncer
date: 2022-03-13 17:54:17
tags:
---

# tx_noncer

> 本部分较为简单，但是还是按照流程进行一些必要的函数分析

### 数据结构

代码文件中的`txNoncer`的数据结构如下：

```go
type txNoncer struct {
	fallback *state.StateDB//对数据库的复制
	nonces   map[common.Address]uint64
	//互斥锁 进行加锁操作 保证安全
	lock     sync.Mutex
}
```

下面是必要的初始化环节：

```go
// newTxNoncer creates a new virtual state database to track the pool nonces.
func newTxNoncer(statedb *state.StateDB) *txNoncer {
	return &txNoncer{
		fallback: statedb.Copy(),
		nonces:   make(map[common.Address]uint64),
	}
}
```

这一部分唯一需要注意的是：`fallback`是直接用`copy()`函数进行复制，而nonces直接初始化为空`map`

---

### 函数解析

1. `get(addr common.Address) uint64`存在该地址对应的`nonce`就直接返回，没有就用该地址对应的`nonce`进行复制，再返回该`nonce值`
2. `set(addr common.Address, nonce uint64)`根据所给的`nonce值`进行赋值即可
3. `setIfLower(addr common.Address, nonce uint64)`将本地址对应的`nonce`与给的`nonce`进行对比，更新为两者之间最小的那个`nonce`，后面的`nonce`对应的交易会进行相应的处理
4. `setAll(all map[common.Address]uint64)`直接将`map`用给定全部复制
