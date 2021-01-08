package business
import (
	"chaincode/nvxclouds_chaincode/entity"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)


//DataUserInvoke 数据节点注册
func DataUserInvoke(stub shim.ChaincodeStubInterface, args []string ) (p pb.Response) {
	defer func() {
		if err := recover(); err != nil {
			p = entity.CatchErr("DataUserInvoke", err)
		}
	}()

	parameter, getParaErr := entity.FilterParameter(stub, args,"invoke")
	if getParaErr != nil {
		return shim.Error(getParaErr.GetJSON())
	}

	_, err := entity.DataUserSave(stub, parameter)
	if err != nil {
		return shim.Error(err.GetJSON())
	}

	resString := "{\"TransactionID\":\"" + stub.GetTxID() + "\"}"
	logger.Info("DataUserInvoke return:", resString)
	return shim.Success([]byte(resString))
}


//DataUserQuery 查询数据使用者结构体信息
func DataUserQuery(stub shim.ChaincodeStubInterface, args []string) (p pb.Response) {
	defer func() {
		if err := recover(); err != nil {
			p = entity.CatchErr("QueryDataUserStruct", err)
		}
	}()

	parameter, getParaErr := entity.FilterParameter(stub, args,"query")
	if getParaErr != nil {
		return shim.Error(getParaErr.GetJSON())
	}

	bytes, ccErr := entity.QueryDataUserStruct(stub, parameter)
	if ccErr != nil {
		logger.Warning(ccErr.GetJSON())
		return shim.Error(ccErr.GetJSON())
	}

	return shim.Success(bytes)
}