pragma solidity ^0.6.0;

/**
 * @title CheckpointOracle
 * @author Gary Rong<garyrong@ethereum.org>, Martin Swende <martin.swende@ethereum.org>
 * @dev Implementation of the blockchain checkpoint registrar.
 */
contract CheckpointOracle {
    /*
        Events
    */

    //生成新的检查点时，有人成功签名

    // NewCheckpointVote is emitted when a new checkpoint proposal receives a vote.
    event NewCheckpointVote(
        uint64 indexed index,
        bytes32 checkpointHash,
        uint8 v,
        bytes32 r,
        bytes32 s
    );

    //构造检查点，管理员地址（需要管理员签名才可以生成检查点）、每一段的大小（每个这么一段的区块数）、需要确认的区块数、阈值
    /*
        Public Functions
    */
    constructor(
        address[] memory _adminlist,
        uint256 _sectionSize,
        uint256 _processConfirms,
        uint256 _threshold
    ) public {
        for (uint256 i = 0; i < _adminlist.length; i++) {
            admins[_adminlist[i]] = true;
            adminList.push(_adminlist[i]);
        }
        sectionSize = _sectionSize;
        processConfirms = _processConfirms;
        threshold = _threshold;
    }

    /**
     * @dev Get latest stable checkpoint information.
     * @return section index
     * @return checkpoint hash
     * @return block height associated with checkpoint
     */
    function GetLatestCheckpoint()
        public
        view
        returns (
            uint64,
            bytes32,
            uint256
        )
    {
        return (sectionIndex, hash, height);
    }

    //注册检查点，最近的区块号和它的哈希值（用于防止链分叉时的重放攻击（第三方拿之前的凭证，谎称自己有签名，冒充身份））、段的索引、段的哈希、签名

    // SetCheckpoint sets  a new checkpoint. It accepts a list of signatures
    // @_recentNumber: a recent blocknumber, for replay protection
    // @_recentHash : the hash of `_recentNumber`
    // @_hash : the hash to set at _sectionIndex
    // @_sectionIndex : the section index to set
    // @v : the list of v-values
    // @r : the list or r-values
    // @s : the list of s-values
    function SetCheckpoint(
        uint256 _recentNumber,
        bytes32 _recentHash,
        bytes32 _hash,
        uint64 _sectionIndex,
        uint8[] memory v,
        bytes32[] memory r,
        bytes32[] memory s
    ) public returns (bool) {
        //已授权

        // Ensure the sender is authorized.
        require(admins[msg.sender]);

        //分叉时，之前的区块位置可能被顶替，区块号相同但是区块哈希不同

        // These checks replay protection, so it cannot be replayed on forks,
        // accidentally or intentionally
        require(blockhash(_recentNumber) == _recentHash);

        //通过签名长度检查签名是否有效

        // Ensure the batch of signatures are valid.
        require(v.length == r.length);
        require(v.length == s.length);

        //还没到下一段，不用新建检查点

        // Filter out "future" checkpoint.
        if (
            block.number < (_sectionIndex + 1) * sectionSize + processConfirms
        ) {
            return false;
        }

        //这一段已经生成过了，错误。

        // Filter out "old" announcement
        if (_sectionIndex < sectionIndex) {
            return false;
        }

        //这一段已经开始生产或者已经生成，没必要再次尝试创建

        // Filter out "stale" announcement
        if (
            _sectionIndex == sectionIndex && (_sectionIndex != 0 || height != 0)
        ) {
            return false;
        }

        //检查点异常，哈希无效
        // Filter out "invalid" announcement
        if (_hash == "") {
            return false;
        }

        // EIP 191 style signatures
        //
        // Arguments when calculating hash to validate
        // 1: byte(0x19) - the initial 0x19 byte
        // 2: byte(0) - the version byte (data with intended validator)
        // 3: this - the validator address
        // --  Application specific data
        // 4 : checkpoint section_index(uint64)
        // 5 : checkpoint hash (bytes32)
        //     hash = keccak256(checkpoint_index, section_head, cht_root, bloom_root)
        bytes32 signedHash = keccak256(
            abi.encodePacked(
                bytes1(0x19),
                bytes1(0),
                this,
                _sectionIndex,
                _hash
            )
        );

        address lastVoter = address(0);

        // In order for us not to have to maintain a mapping of who has already
        // voted, and we don't want to count a vote twice, the signatures must
        // be submitted in strict ordering.
        for (uint256 idx = 0; idx < v.length; idx++) {
            address signer = ecrecover(signedHash, v[idx], r[idx], s[idx]);
            require(admins[signer]);
            require(uint256(signer) > uint256(lastVoter));
            lastVoter = signer;
            emit NewCheckpointVote(
                _sectionIndex,
                _hash,
                v[idx],
                r[idx],
                s[idx]
            );

            // Sufficient signatures present, update latest checkpoint.
            if (idx + 1 >= threshold) {
                hash = _hash;
                height = block.number;
                sectionIndex = _sectionIndex;
                return true;
            }
        }
        // We shouldn't wind up here, reverting un-emits the events
        revert();
    }

    /**
     * @dev Get all admin addresses
     * @return address list
     */
    function GetAllAdmin() public view returns (address[] memory) {
        address[] memory ret = new address[](adminList.length);
        for (uint256 i = 0; i < adminList.length; i++) {
            ret[i] = adminList[i];
        }
        return ret;
    }

    /*
        Fields
    */
    // A map of admin users who have the permission to update CHT and bloom Trie root
    mapping(address => bool) admins;

    // A list of admin users so that we can obtain all admin users.
    address[] adminList;

    // Latest stored section id
    uint64 sectionIndex;

    // The block height associated with latest registered checkpoint.
    uint256 height;

    // The hash of latest registered checkpoint.
    bytes32 hash;

    // The frequency for creating a checkpoint
    //
    // The default value should be the same as the checkpoint size(32768) in the ethereum.
    uint256 sectionSize;

    // The number of confirmations needed before a checkpoint can be registered.
    // We have to make sure the checkpoint registered will not be invalid due to
    // chain reorg.
    //
    // The default value should be the same as the checkpoint process confirmations(256)
    // in the ethereum.
    uint256 processConfirms;

    // The required signatures to finalize a stable checkpoint.
    uint256 threshold;
}
