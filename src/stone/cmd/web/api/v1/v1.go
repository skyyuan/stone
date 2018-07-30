package v1

import (
	"strconv"

	ecom "github.com/ethereum/go-ethereum/common"
	"github.com/labstack/echo"
	"stone/common"
	"stone/service/web"
)

// RegisterAPI 注册v1版本的API
func RegisterAPI(router *echo.Echo) {
	v1 := router.Group("/v1")

	v1.POST("/get_erc20cointypelist", func(c echo.Context) error {
		erc20s := web.GetErc20CoinTypeList()
		return common.JSONReturns(c, erc20s)
	})

	v1.POST("/get_estimategas", func(c echo.Context) error {
		var params EstimateGasPayload
		if err := c.Bind(&params); err != nil {
			return err
		}

		if params.ContractAddress == "" || ecom.IsHexAddress(params.ContractAddress) {
			estimateGas := web.GetEstimateGas(params.ContractAddress)
			estimateGasRly := &EthereumEstimateGasResponse{
				ContractAddress: estimateGas.ContractAddress,
				EstimateGas:     strconv.FormatInt(estimateGas.EstimateGas, 10),
			}
			return common.JSONReturns(c, estimateGasRly)
		}
		return common.BizError1001
	})

	v1.POST("/get_allestimategases", func(c echo.Context) error {
		var estimateGases []*EthereumEstimateGasResponse

		for _, estimateGas := range web.GetAllEstimateGases() {
			estimateGases = append(estimateGases, &EthereumEstimateGasResponse{
				ContractAddress: estimateGas.ContractAddress,
				EstimateGas:     strconv.FormatInt(estimateGas.EstimateGas, 10),
			})
		}

		return common.JSONReturns(c, estimateGases)
	})
	v1.POST("/register_address", func(c echo.Context) error {
		var params RegisterAddress
		if err := c.Bind(&params); err != nil {
			return err
		}
		if ecom.IsHexAddress(params.Address) {
			err := web.RegisterAddress(params.Address, params.AddressType, params.DeviceType, params.DeviceUUID)
			if err != nil {
				return err
			}
			return common.JSONReturns(c, nil)
		}
		return common.BizError1001
	})
	v1.POST("/user_erc20_tokens", func(c echo.Context) error {
		var params EthAddress
		if err := c.Bind(&params); err != nil {
			return err
		}
		if ecom.IsHexAddress(params.Address) {
			erc20Coins, err := web.UserErc20TokenList(params.Address)
			if err != nil {
				return err
			}
			return common.JSONReturns(c, erc20Coins)
		}
		return common.BizError1001
	})
	v1.POST("/fetch_new_user_erc20_tokens", func(c echo.Context) error {
		var params EthAddresses
		res := make(map[string]([]web.Erc20CoinType))
		if err := c.Bind(&params); err != nil {
			return err
		}
		for _, item := range params.Addresses {
			if !ecom.IsHexAddress(item) {
				return common.BizError1001
			}
		}
		for _, item := range params.Addresses {
			erc20Coins, err := web.FetchNewUserErc20TokenList(item)
			if err != nil {
				return err
			}
			res[item] = erc20Coins
		}
		return common.JSONReturns(c, res)
	})
}
