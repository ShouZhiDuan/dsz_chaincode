package business

import (
	"chaincode/nvxclouds_chaincode/entity"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

//AppUserInvoke 数据节点注册
func AppUserInvoke(stub shim.ChaincodeStubInterface, args []string) (p pb.Response) {
	defer func() {
		if err := recover(); err != nil {
			p = entity.CatchErr("AppUserInvoke", err)
		}
	}()

	parameter, getParaErr := entity.FilterParameter(stub, args, "invoke")
	if getParaErr != nil {
		return shim.Error(getParaErr.GetJSON())
	}

	_, err := entity.AppUserSave(stub, parameter)
	if err != nil {
		return shim.Error(err.GetJSON())
	}

	resString := "{\"TransactionID\":\"" + stub.GetTxID() + "\"}"
	logger.Info("AppUserInvoke return:", resString)
	return shim.Success([]byte(resString))
}

//AppUserQuery 查询应用使用者结构体信息
func AppUserQuery(stub shim.ChaincodeStubInterface, args []string) (p pb.Response) {
	defer func() {
		if err := recover(); err != nil {
			p = entity.CatchErr("AppUserQuery", err)
		}
	}()

	parameter, getParaErr := entity.FilterParameter(stub, args, "query")
	if getParaErr != nil {
		return shim.Error(getParaErr.GetJSON())
	}

	bytes, ccErr := entity.QueryAppUserStruct(stub, parameter)
	if ccErr != nil {
		logger.Warning(ccErr.GetJSON())
		return shim.Error(ccErr.GetJSON())
	}

	return shim.Success(bytes)
}

//AppUserQueryBalance 查询积分余额
func AppUserQueryBalance(stub shim.ChaincodeStubInterface, args []string) (p pb.Response) {
	defer func() {
		if err := recover(); err != nil {
			p = entity.CatchErr("business AppUserQueryBalance", err)
		}
	}()

	parameter, getParaErr := entity.FilterParameter(stub, args, "query")
	if getParaErr != nil {
		return shim.Error(getParaErr.GetJSON())
	}

	var adubOperation, aapirOperation []entity.AppUserStruct
	parameter.Operation = "authorizationDataUpBlockChain"
	bytes, _ := entity.QueryAppUserStruct(stub, parameter)
	if bytes != nil {
		errs := json.Unmarshal(bytes, &adubOperation)
		if errs != nil {
			ccErr := entity.GetError(1006)
			ccErr.AddMessage("business AppUserQueryBalance get:" + errs.Error())
			return shim.Error(ccErr.GetJSON())
		}
	}

	parameter.Operation = "authorizationAndParticipationInResearch"
	bytess, _ := entity.QueryAppUserStruct(stub, parameter)
	if bytess != nil {
		errs := json.Unmarshal(bytess, &aapirOperation)
		if errs != nil {
			ccErr := entity.GetError(1006)
			ccErr.AddMessage("business AppUserQueryBalance get:" + errs.Error())
			return shim.Error(ccErr.GetJSON())
		}
	}

	var balance int64
	if adubOperation != nil {
		for _, val := range adubOperation {
			balance += val.Balance
		}
	}

	if aapirOperation != nil {
		for _, val := range aapirOperation {
			balance += val.Balance
		}
	}

	balanceStruct := entity.BalanceStruct{
		AppUserID: parameter.AppUserID,
		Balance:   balance,
	}
	balanceByte, err := json.Marshal(balanceStruct)
	if err != nil {
		ccErr := entity.GetError(1005)
		ccErr.AddMessage("business AppUserQueryBalance :" + err.Error())
		logger.Warning(ccErr.GetJSON())
		return shim.Error(ccErr.GetJSON())
	}

	return shim.Success(balanceByte)
}

//AppUserQueryBalanceHistory 查询积分历史
func AppUserQueryBalanceHistory(stub shim.ChaincodeStubInterface, args []string) (p pb.Response) {
	defer func() {
		if err := recover(); err != nil {
			p = entity.CatchErr("business AppUserQueryBalance", err)
		}
	}()

	parameter, getParaErr := entity.FilterParameter(stub, args, "query")
	if getParaErr != nil {
		return shim.Error(getParaErr.GetJSON())
	}

	var adubOperation, aapirOperation, returnHistory []entity.AppUserStruct
	parameter.Operation = "authorizationDataUpBlockChain"
	bytes, _ := entity.QueryAppUserStruct(stub, parameter)
	if bytes != nil {
		errs := json.Unmarshal(bytes, &adubOperation)
		if errs != nil {
			ccErr := entity.GetError(1006)
			ccErr.AddMessage("business AppUserQueryBalance get:" + errs.Error())
			return shim.Error(ccErr.GetJSON())
		}
	}

	parameter.Operation = "authorizationAndParticipationInResearch"
	bytess, _ := entity.QueryAppUserStruct(stub, parameter)
	if bytess != nil {
		errs := json.Unmarshal(bytess, &aapirOperation)
		if errs != nil {
			ccErr := entity.GetError(1006)
			ccErr.AddMessage("business AppUserQueryBalance get:" + errs.Error())
			return shim.Error(ccErr.GetJSON())
		}
	}

	if adubOperation != nil {
		for _, val := range adubOperation {
			returnHistory = append(returnHistory, val)
		}
	}

	if aapirOperation != nil {
		for _, val := range aapirOperation {
			returnHistory = append(returnHistory, val)
		}
	}

	//按照时间先后顺序排序
	history := entity.SortListByTime(returnHistory)

	balanceByte, err := json.Marshal(history)
	if err != nil {
		ccErr := entity.GetError(1005)
		ccErr.AddMessage("business AppUserQueryBalance :" + err.Error())
		logger.Warning(ccErr.GetJSON())
		return shim.Error(ccErr.GetJSON())
	}

	return shim.Success(balanceByte)
}
