package business
import (
	"chaincode/nvxclouds_chaincode/entity"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)


//RegulatorInvoke 数据节点注册
func RegulatorInvoke(stub shim.ChaincodeStubInterface, args []string ) (p pb.Response) {
	defer func() {
		if err := recover(); err != nil {
			p = entity.CatchErr("RegulatorInvoke", err)
		}
	}()

	parameter, getParaErr := entity.FilterParameter(stub, args,"invoke")
	if getParaErr != nil {
		return shim.Error(getParaErr.GetJSON())
	}

	_, err := entity.RegulatorSave(stub, parameter)
	if err != nil {
		return shim.Error(err.GetJSON())
	}

	resString := "{\"TransactionID\":\"" + stub.GetTxID() + "\"}"
	logger.Info("RegulatorInvoke return:", resString)
	return shim.Success([]byte(resString))
}

//RegulatorQuery 查询监管方结构体信息
func RegulatorQuery(stub shim.ChaincodeStubInterface, args []string) (p pb.Response) {
	defer func() {
		if err := recover(); err != nil {
			p = entity.CatchErr("QueryRegulatorStruct", err)
		}
	}()

	parameter, getParaErr := entity.FilterParameter(stub, args,"query")
	if getParaErr != nil {
		return shim.Error(getParaErr.GetJSON())
	}

	bytes, ccErr := entity.QueryRegulatorStruct(stub, parameter)
	if ccErr != nil {
		logger.Warning(ccErr.GetJSON())
		return shim.Error(ccErr.GetJSON())
	}

	return shim.Success(bytes)
}