package config

import (
	"github.com/ethereum/go-ethereum/common"
)

var (
	// chain ID
	goerliChainID = int64(5)
	xdaiChainID   = int64(100)

	// start block
	goerliStartBlock = uint64(5184980)
	xdaiStartBlock   = uint64(16515648)

	// factory address
	goerliContractAddress            = common.HexToAddress("0x2bABAbaCcb615275E68168FB51609Ff08E21889A")
	goerliFactoryAddress             = common.HexToAddress("0x539F6636e068d194E2b8737C4d66fD2A89e73657")
	goerliPriceOracleContractAddress = common.HexToAddress("0xBcD55D9B095E91FA84A127cf36a2039Ed7896366")

	xdaiContractAddress = common.HexToAddress("0x0FDc5429C50e2a39066D8A94F3e2D2476fcc3b85")
	xdaiFactoryAddress  = common.HexToAddress("0xc2d5a532cf69aa9a1378737d8ccdef884b6e7420")

	// postage stamp
	goerliPostageStampContractAddress = common.HexToAddress("0x621e455C4a139f5C4e4A8122Ce55Dc21630769E4")
	xdaiPostageStampContractAddress   = common.HexToAddress("0x6a1a21eca3ab28be85c7ba22b2d6eae5907c900e")
)

type ChainConfig struct {
	StartBlock         uint64
	LegacyFactories    []common.Address
	PostageStamp       common.Address
	CurrentFactory     common.Address
	PriceOracleAddress common.Address
}

func GetChainConfig(chainID int64) (*ChainConfig, bool) {
	var cfg ChainConfig
	switch chainID {
	case goerliChainID:
		cfg.PostageStamp = goerliPostageStampContractAddress
		cfg.StartBlock = goerliStartBlock
		cfg.CurrentFactory = goerliFactoryAddress
		cfg.LegacyFactories = []common.Address{}
		cfg.PriceOracleAddress = goerliPriceOracleContractAddress
		return &cfg, true
	case xdaiChainID:
		cfg.PostageStamp = xdaiPostageStampContractAddress
		cfg.StartBlock = xdaiStartBlock
		cfg.CurrentFactory = xdaiFactoryAddress
		cfg.LegacyFactories = []common.Address{}
		cfg.PriceOracleAddress = xdaiContractAddress
		return &cfg, true
	default:
		return &cfg, false
	}
}
