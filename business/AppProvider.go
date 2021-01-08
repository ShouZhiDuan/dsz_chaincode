package business

import (
	"chaincode/nvxclouds_chaincode/entity"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

//AppProviderInvoke 应用提供方写入操作
func AppProviderInvoke(stub shim.ChaincodeStubInterface, args []string) (p pb.Response) {
	defer func() {
		if err := recover(); err != nil {
			p = entity.CatchErr("AppProviderInvoke", err)
		}
	}()

	parameter, getParaErr := entity.FilterParameter(stub, args,"invoke")
	if getParaErr != nil {
		return shim.Error(getParaErr.GetJSON())
	}

	_, err := entity.AppProviderSave(stub, parameter)
	if err != nil {
		return shim.Error(err.GetJSON())
	}

	resString := "{\"TransactionID\":\"" + stub.GetTxID() + "\"}"
	logger.Info("AppProviderInvoke return:", resString)
	return shim.Success([]byte(resString))
}

//AppProviderQuery 查询应用提供方结构体信息
func AppProviderQuery(stub shim.ChaincodeStubInterface, args []string) (p pb.Response) {
	defer func() {
		if err := recover(); err != nil {
			p = entity.CatchErr("AppProviderQuery", err)
		}
	}()

	parameter, getParaErr := entity.FilterParameter(stub, args,"query")
	if getParaErr != nil {
		return shim.Error(getParaErr.GetJSON())
	}

	bytes, ccErr := entity.QueryAppProviderStruct(stub, parameter)
	if ccErr != nil {
		logger.Warning(ccErr.GetJSON())
		return shim.Error(ccErr.GetJSON())
	}

	return shim.Success(bytes)
}
