package v1

// EstimateGasPayload represents payload for getting estimateGas
type EstimateGasPayload struct {
	ContractAddress string `json:"contract_address"`
}

// EthereumEstimateGasResponse represents estimate Gas used when make a call or tx
type EthereumEstimateGasResponse struct {
	// ContractAddress is "" when ETH
	ContractAddress string `json:"contractAddress"`
	EstimateGas     string `json:"estimateGas"`
}

// RegisterAddress RegisterAddress payload
type RegisterAddress struct {
	Address     string `json:"address"`
	AddressType uint   `json:"address_type"`
	DeviceType  uint   `json:"device_type"`
	DeviceUUID  string `json:"device_uuid"`
}

type EthAddress struct {
	Address string `json:"address"`
}

// EthAddresses EthAddresses
type EthAddresses struct {
	Addresses []string `json:"addresses"`
}
