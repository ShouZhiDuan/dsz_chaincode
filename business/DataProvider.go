package business

import (
	"chaincode/nvxclouds_chaincode/entity"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("business")

//DataProviderInvoke 数据节点注册
func DataProviderInvoke(stub shim.ChaincodeStubInterface, args []string) (p pb.Response) {
	defer func() {
		if err := recover(); err != nil {
			p = entity.CatchErr("DataProviderInvoke", err)
		}
	}()

	parameter, getParaErr := entity.FilterParameter(stub, args,"invoke")
	if getParaErr != nil {
		return shim.Error(getParaErr.GetJSON())
	}

	_, err := entity.DataProviderSave(stub, parameter)
	if err != nil {
		return shim.Error(err.GetJSON())
	}

	resString := "{\"TransactionID\":\"" + stub.GetTxID() + "\"}"
	logger.Info("DataProviderInvoke return:", resString)
	return shim.Success([]byte(resString))
}

//DataProviderQuery 查询数据提供方结构体信息
func DataProviderQuery(stub shim.ChaincodeStubInterface, args []string) (p pb.Response) {
	defer func() {
		if err := recover(); err != nil {
			p = entity.CatchErr("QueryDataProviderStruct", err)
		}
	}()

	parameter, getParaErr := entity.FilterParameter(stub, args,"query")
	if getParaErr != nil {
		return shim.Error(getParaErr.GetJSON())
	}

	bytes, ccErr := entity.QueryDataProviderStruct(stub, parameter)
	if ccErr != nil {
		logger.Warning(ccErr.GetJSON())
		return shim.Error(ccErr.GetJSON())
	}

	return shim.Success(bytes)
}
