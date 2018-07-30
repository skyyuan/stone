package eth

const (
	// HexPrefix 16进制前缀
	HexPrefix = "0x"

	// ERC20MethodTransfer erc20方法transfer的16进制字符串
	ERC20MethodTransfer = "a9059cbb"

	// ERC20TransferLength erc20方法transfer的input长度
	ERC20TransferLength = 130

	// ERC20MethodBalanceOf erc20方法balanceOf的16进制字符串
	ERC20MethodBalanceOf = "70a08231"

	// ERC20Name erc20 name 的16进制字符串
	ERC20Name = "0x06fdde03"
	// ERC20Symbol erc20 symbol 的16进制字符串
	ERC20Symbol = "0x95d89b41"
	// ERC20Decimals erc20 decimals 的16进制字符串
	ERC20Decimals = "0x313ce567"
	// ERC20AbiDefaultLength ERC20 Abi字符串 Default Length
	ERC20AbiDefaultLength = 194

	// ConfirmedNum confirmed num till to 12
	ConfirmedNum = int64(12)

	// EthGasLimit eth gas limit transaction gas limit
	EthGasLimit = int64(21000)
	// DefaultErc20GasPrice default gas price for erc20
	DefaultErc20GasPrice = int64(60000)

	// DefaultErc20Icon default erc20 icon
	DefaultErc20Icon = "https://ops.58wallet.io/home/img/avatar-DEFAULT@2x.png"

	// StatusPass status pass
	StatusPass = 0
	// StatusFailed status failed
	StatusFailed = 1
	// StatusPending status pending
	StatusPending = 2

	// EtherscanURL EtherscanURL
	EtherscanURL = "https://etherscan.io/tx/"
	// RopstenURL RopstenURL
	RopstenURL = "https://ropsten.etherscan.io/tx/"

	// RopstenTestNode Ropsten test node
	RopstenTestNode = "http://47.52.31.232:8545"

	// EventTransferHash hash of Event Transfer
	EventTransferHash = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
	// AddressInHashIndex Address in hash  from index
	AddressInHashIndex = 26

	// MyEtherWalletEndpointTestNet MyEtherWalletEndpointTestNet
	MyEtherWalletEndpointTestNet = "https://api.myetherapi.com/rop"
	// MyEtherWalletEndpoint MyEtherWalletEndpoint
	MyEtherWalletEndpoint = "https://api.myetherapi.com/eth"
)

var ATMFunction map[string]string = map[string]string{
	"b6665eac": "ATM()",
	"d53d8c1a": "ROLE_ATMPLATFORM()",
	"98a1b397": "ROLE_OPERATOR()",
	"a4114bdc": "ROLE_SCREEN()",
	"b941719e": "ROLE_USER()",
	"79ba5097": "acceptOwnership()",
	"9e8674dc": "atm()",
	"674acbd7": "atmPlatformAddr()",
	"a6f9dae1": "changeOwner(address)",
	"f2ae98b2": "deleteAdvertise(uint256)",
	"93e3168a": "enableIssue(uint256)",
	"2728f93e": "freeMediaIndex(uint256)",
	"99d24319": "getMediaStrategy(uint256)",
	"1f5e05bf": "issueReward(address,address,uint256,uint256)",
	"d44708e0": "mediaList(address,uint256)",
	"f6854c3c": "mediaStrategy(uint256)",
	"dd20dc8c": "pauseIssue(uint256)",
	"5396b906": "publishAdvertise(address,uint256,uint256,uint256,uint256[3])",
	"d87b9ead": "updateAdvertise(address,uint256,uint256,uint256,uint256[3])",
	"32bc934c": "MILLION()",
	"dd62ed3e": "allowance(address,address)",
	"095ea7b3": "approve(address,uint256)",
	"70a08231": "balanceOf(address)",
	"27e235e3": "balances(address)",
	"313ce567": "decimals()",
	"ee070805": "disabled()",
	"681d8345": "getATMTotalSupply()",
	"41c0e1b5": "kill()",
	"06fdde03": "name()",
	"d4ee1d90": "newOwner()",
	"8da5cb5b": "owner()",
	"6c5a7d1e": "setDisabled(bool)",
	"95d89b41": "symbol()",
	"18160ddd": "totalSupply()",
	"a9059cbb": "transfer(address,uint256)",
	"23b872dd": "transferFrom(address,address,uint256)",
	"54fd4d50": "version()",
}
