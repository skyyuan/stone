package v1

import "stone/common"

// EthereumTransactionListPayload request payload
type EthereumTransactionListPayload struct {
	UserAddress     string `json:"user_address"`
	ContractAddress string `json:"contract_address"`
	common.PageParams
}

// EthereumAccount represents an ethereum account
type EthereumAccount struct {
	UserAddress     string `json:"user_address"`
	ContractAddress string `json:"contract_address"`
}

// EthereumTransactionHash represents an ethereum transaction hash
type EthereumTransactionHash struct {
	TxHash string `json:"tx_hash"`
}

// EthereumAccountMultipleContract represents an account and multiple contracts to which the account belong.
type EthereumAccountMultipleContract struct {
	UserAddress       string   `json:"user_address"`
	ContractAddresses []string `json:"contract_addresses"`
}

// EthereumSendRawTransactionPayload represents payload of request to send a raw tx
type EthereumSendRawTransactionPayload struct {
	SignedTransaction string `json:"signed_transaction"`
}

// EthereumSendRawTransaction represents response for a raw tx
type EthereumSendRawTransaction struct {
	Hash string `json:"hash"`
}

// EthereumTxnCountPayload get a useraddress txn count
type EthereumTxnCountPayload struct {
	UserAddress string `json:"user_address"`
}

// EthereumTransactionListPayload request payload
type EthereumBlockPayload struct {
	BlockNumber string `json:"block_number"`
	common.PageParams
}

// EthereumTransactionListPayload request payload
type EthereumBlockListPayload struct {
	common.PageParams
}
type EthereumATMpricePayload struct {
	common.PageParams
}
