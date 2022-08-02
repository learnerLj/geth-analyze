## 特殊函数

### 底层函数交互

**特殊交互方式**

call 是底层的调用（没有封装过），直接发送消息给合约。方式如下：

1. 所有的参数，都会打包成一串32个字节，连续存放的序列。
2. 若第一个参数是函数的签名（函数哈希之后的前4个字节），则第二、第三这些后面的参数是函数的参数。如：`nameReg.call(bytes4(keccak256("fun(uint256)")), a);`

- `<address>.call(bytes memory) returns (bool, bytes memory)`

  用给定的合约发出低级 `CALL` 调用，返回成功状态及返回数据，发送所有可用 gas，也可以调节 gas。

```js
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.10;

contract Receiver {
    event Received(address caller, uint amount, string message);

    fallback() external payable {//回退函数
        emit Received(msg.sender, msg.value, "Fallback was called");
    }

    function foo(string memory _message, uint _x) public payable returns (uint) {
        emit Received(msg.sender, msg.value, _message);

        return _x + 1;
    }
}

contract Caller {
    event Response(bool success, bytes data);

    function testCallFoo(address payable _addr) public payable {
        // 注意观察调用的格式，结合前面学习的 ABI 编码方式。
        (bool success, bytes memory data) = _addr.call{value: msg.value, gas: 5000}(
            abi.encodeWithSignature("foo(string,uint256)", "call foo", 123)
        );

        emit Response(success, data);
    }

    // 不存在的函数调用会失败，但是同样会触发回调函数。
    function testCallDoesNotExist(address _addr) public {
        (bool success, bytes memory data) = _addr.call(
            abi.encodeWithSignature("doesNotExist()")
        );

        emit Response(success, data);
    }
}

```



- `<address>.delegatecall(bytes memory) returns (bool, bytes memory)`

  用给定的合约发出低级 `DELEGATECALL` 调用 ，返回成功状态并返回数据，失败时返回 `false`。上下文属于发出调用的合约。

```js
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.10;

// NOTE:先部署这个合约
contract B {
    // NOTE: storage layout must be the same as contract A
    uint public num;
    address public sender;
    uint public value;

    function setVars(uint _num) public payable {
        num = _num;
        sender = msg.sender;
        value = msg.value;
    }
}

contract A {
    uint public num;
    address public sender;
    uint public value;

    function setVars(address _contract, uint _num) public payable {
        // 只改变了合约A的值，因为上下文属于合约A。
        (bool success, bytes memory data) = _contract.delegatecall(
            abi.encodeWithSignature("setVars(uint256)", _num)
        );
    }
}

```

- `<address>.staticcall(bytes memory) returns (bool, bytes memory)`

  用给定的有效载荷 发出低级 `STATICCALL` 调用 ，如果改变了被调用合约的状态，立即回滚。

注意：`.call`会绕过类型检查，函数存在检查和参数打包。

注意：`send`调用栈深度达到1024就会失败。

注意：0.5.0以后不允许通过合约实例来访问地址成员`this.balance`。0.5.0以前，底层调用只会返回是否成功不会返回数据。

注意：因为EVM不会检查调用的合约是否存在，并且总是把调用不存在的合约视为成功，因此提供了 `extcodesize` 的操作码，确认合约存在（即合约地址内有代码），否则引起异常。注意底层调用不会触发。

注意：底层的调用略去了很多检查，使得他们更加节省gas但是安全性更低。



### receive 函数

一个合约至多有一个 receive 函数，形如 ``receive() external payable { ... }`` ，注意：

- 没有 `function` 的标识
- 没有参数
- 只能是 `external` 和`payable`标识
- 可以有函数修饰器
- 支持重载。

**`receive` 函数在调用数据为空时（如用 `call` 传入空字节，或者转账）执行**,如果没有设置 `receive` 函数，那么就会执行 `fallback` 函数，如果这两个函数都不存在，合约就不能通过交易的形式获取以太币。

**注意 `receive`函数只有2300gas可用，因此它进行其他操作的空间很小。** 以下功能都因为超过消耗的 gas 而不能够实现。

- 写入存储
- 创建合约
- 调用消耗大量 gas 的外部函数
- 发送以太币

每一步都会消耗2300gas.

我们建议**只使用 `receive` 函数来接收以太币**。

### 回退函数

一个合约至多一个回退函数，格式如：`fallback () external [payable]` 或者 `fallback (bytes calldata _input) external [payable] returns (bytes memory _output)`，后者的函数参数会接收完整的调用数据（`msg.data`)，返回未经过ABI编码的原始数据。

- 回退函数只当没有与调用数据匹配的函数签名时执行。
- 可以重载，也可以被修饰器修饰。
- 在函数调用时，如果**没有与之匹配的函数签名或者调用数据为空且无`receive`函数**，就会调用 `fallbakc` 函数。
- 如果回退函数代替了`receive`函数完成接收以太币的功能，那么仍然只有2300gas可用。