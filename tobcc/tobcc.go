package main

import (
	"chaincode/nvxclouds_chaincode/business"
	"chaincode/nvxclouds_chaincode/entity"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("TOBCC")

func main() {
	err := shim.Start(new(TOBCC))
	if err != nil {
		logger.Error("Error starting TOBCC chaincode: %s", err)
	}
}

//TOBCC chaincode名称
type TOBCC struct{}

//Init 初始化
func (cc *TOBCC) Init(stub shim.ChaincodeStubInterface) (p pb.Response) {
	defer func() {
		if err := recover(); err != nil {
			p = entity.CatchErr("Init", err)
		}
	}()
	logger.Info("[TOBCC][Init]...")
	return shim.Success(nil)
}

//Invoke :
func (cc *TOBCC) Invoke(stub shim.ChaincodeStubInterface) (p pb.Response) {
	defer func() {
		if err := recover(); err != nil {
			p = entity.CatchErr("Invoke", err)
		}
	}()
	fmt.Println("TOBCC username: ====", entity.GetUserName(stub))

	function, args := stub.GetFunctionAndParameters()
	logger.Info("[TOBCC] function:", function, ";  args:", args)
	//-------------------------数据提供方-----------------------------------
	if function == "dataProviderInvoke" { //数据提供方写入操作
		return business.DataProviderInvoke(stub, args)
	} else if function == "dataProviderQuery" { //查询数据提供方数据
		return business.DataProviderQuery(stub, args)
		//-------------------------应用提供方-----------------------------------
	} else if function == "appProviderInvoke" { //应用提供方写入操作
		return business.AppProviderInvoke(stub, args)
	} else if function == "appProviderQuery" { //查询应用提供方数据
		return business.AppProviderQuery(stub, args)
		//-------------------------数据使用方-----------------------------------
	} else if function == "dataUserInvoke" { //数据使用方写入操作
		return business.DataUserInvoke(stub, args)
	} else if function == "dataUserQuery" { //查询数据使用方数据
		return business.DataUserQuery(stub, args)
		//-------------------------监管方-----------------------------------
	} else if function == "regulatorInvoke" { //监管方写入操作
		return business.RegulatorInvoke(stub, args)
	} else if function == "regulatorQuery" { //查询监管方数据
		return business.RegulatorQuery(stub, args)
	} else {
		nvxcloudsErr := entity.GetError(1301)
		errMessage := entity.GetErrMessageByIndentity(entity.GetUserName(stub), function)
		nvxcloudsErr.AddMessage(errMessage)
		return shim.Error(nvxcloudsErr.GetJSON())
	}
}
