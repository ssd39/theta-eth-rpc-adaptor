package ethrpc

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/thetatoken/theta-eth-rpc-adaptor/common"
	rpcc "github.com/ybbus/jsonrpc"

	tcommon "github.com/thetatoken/theta/common"
	"github.com/thetatoken/theta/ledger/types"

	// "github.com/thetatoken/theta/ledger/types"
	trpc "github.com/thetatoken/theta/rpc"
)

type TxTmp struct {
	Tx   json.RawMessage `json:"raw"`
	Type byte            `json:"type"`
	Hash tcommon.Hash    `json:"hash"`
}

// ------------------------------- eth_gasPrice -----------------------------------

func (e *EthRPCService) GasPrice(ctx context.Context) (result string, err error) {
	logger.Infof("eth_gasPrice called")
	gasPrice := big.NewInt(4100)
	fmt.Printf("gasPrice: %v\n", gasPrice)
	result = "0x" + gasPrice.Text(16)
	return result, nil
}

func getDefaultGasPrice(client *rpcc.RPCClient) *big.Int {
	gasPrice := big.NewInt(0) // Default for the Main Chain
	ethChainID, err := getEthChainID(client)
	if err == nil {
		if ethChainID > 1000 { // must be a Subchain
			gasPrice = big.NewInt(1000) // Default for the Subchains
		}
	}
	return gasPrice
}

func getEthChainID(client *rpcc.RPCClient) (uint64, error) {
	rpcRes, rpcErr := client.Call("theta.GetStatus", trpc.GetStatusArgs{})
	var blockHeight uint64
	parse := func(jsonBytes []byte) (interface{}, error) {
		trpcResult := trpc.GetStatusResult{}
		json.Unmarshal(jsonBytes, &trpcResult)
		re := chainIDResultWrapper{
			chainID: trpcResult.ChainID,
		}
		blockHeight = uint64(trpcResult.LatestFinalizedBlockHeight)
		return re, nil
	}

	resultIntf, err := common.HandleThetaRPCResponse(rpcRes, rpcErr, parse)
	if err != nil {
		return 0, err
	}
	thetaChainIDResult, ok := resultIntf.(chainIDResultWrapper)
	if !ok {
		return 0, fmt.Errorf("failed to convert chainIDResultWrapper")
	}

	thetaChainID := thetaChainIDResult.chainID
	ethChainID := types.MapChainID(thetaChainID, blockHeight).Uint64()

	return ethChainID, nil
}
