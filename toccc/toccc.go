package main

import (
	"chaincode/nvxclouds_chaincode/business"
	"chaincode/nvxclouds_chaincode/entity"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("TOCCC")

func main() {
	err := shim.Start(new(TOCCC))
	if err != nil {
		logger.Error("Error starting TOCCC chaincode: %s", err)
	}
}

//TOCCC chaincode名称
type TOCCC struct{}

//Init 初始化
func (cc *TOCCC) Init(stub shim.ChaincodeStubInterface) (p pb.Response) {
	defer func() {
		if err := recover(); err != nil {
			p = entity.CatchErr("Init", err)
		}
	}()
	logger.Info("[TOCCC][Init]...")
	return shim.Success(nil)
}

//Invoke :
func (cc *TOCCC) Invoke(stub shim.ChaincodeStubInterface) (p pb.Response) {
	defer func() {
		if err := recover(); err != nil {
			p = entity.CatchErr("Invoke", err)
		}
	}()
	fmt.Println("TOCCC username: ====", entity.GetUserName(stub))

	function, args := stub.GetFunctionAndParameters()
	logger.Info("[TOCCC] function:", function, ";  args:", args)
	//-------------------------应用使用方-----------------------------------
	if function == "appUserInvoke" { //应用使用方写入操作
		return business.AppUserInvoke(stub, args)
	} else if function == "appUserQuery" { //查看应用使用方数据
		return business.AppUserQuery(stub, args)
	} else if function == "appUserQueryBalance" { //查看积分余额
		return business.AppUserQueryBalance(stub, args)
	} else if function == "appUserQueryBalanceHistory" { //查看积分历史
		return business.AppUserQueryBalanceHistory(stub, args)
	} else {
		nvxcloudsErr := entity.GetError(1301)
		errMessage := entity.GetErrMessageByIndentity(entity.GetUserName(stub), function)
		nvxcloudsErr.AddMessage(errMessage)
		return shim.Error(nvxcloudsErr.GetJSON())
	}
}
