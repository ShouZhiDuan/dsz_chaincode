package entity

import (
	"encoding/json"
	"github.com/goinggo/mapstructure"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strconv"
)

//Save :(it *AppUserStruct)保存账户
//====================================================================
//compositeIndexName 复合键
func (it *AppUserStruct) Save(stub shim.ChaincodeStubInterface, compositeIndexName string) ([]byte, *CcError) {
	args, err := it.getAppUserStructPrefix(stub, compositeIndexName)
	if err != nil {
		return nil, err
	}

	bytes, errs := Set(stub, args)
	if errs != nil {
		return nil, errs
	}
	return bytes, nil
}

//getAppUserStructPrefix :(it *AppUserStruct)查询该操作存储在链上的信息前缀
//====================================================================
func (it *AppUserStruct) getAppUserStructPrefix(stub shim.ChaincodeStubInterface, compositeIndexName string) ([]string, *CcError) {
	appUserStruct := Struct2Map(*it)
	appUserStructBytes, err := json.Marshal(appUserStruct)
	if err != nil {
		ccErr := GetError(1006)
		ccErr.AddMessage("entity Save:" + err.Error())
		return nil, &ccErr
	}
	appUserStructBytesStr := string(appUserStructBytes)

	balanceStr := strconv.FormatInt(it.Balance, 10)
	
	var args []string
	args = []string{compositeIndexName}
	switch compositeIndexName {
	case "authorizationDataUpBlockChain":
		args = append(args, it.AppUserID, balanceStr, it.TxTime, it.MessageHash)
	case "authorizationAndParticipationInResearch":
		args = append(args, it.AppUserID, it.ResearchProjectID, balanceStr, it.TxTime, it.MessageHash)
	case "queryDataUsage":
		args = append(args, it.AppUserID, it.TxTime, it.MessageHash)
	case "queryAppUserDataResearchResults":
		args = append(args, it.AppUserID, it.ResearchProjectID, it.TxTime, it.MessageHash)
	default:
		ccErr := GetError(1301)
		ccErr.AddMessage("entity getAppUserStructPrefix")
		return nil, &ccErr
	}

	args = append(args, appUserStructBytesStr)
	return args, nil
}

//AppUserSave ：应用使用者数据写入
// ========================================================================
func AppUserSave(stub shim.ChaincodeStubInterface, parameter Parameter) ([]byte, *CcError) {
	var balanceInt int64
	balanceInt = 0
	if parameter.Balance == "" {
		if parameter.Operation == "authorizationDataUpBlockChain" || parameter.Operation == "authorizationAndParticipationInResearch" {
			ccErr := GetError(1002)
			ccErr.AddMessage("Balance is empty")
			return nil, &ccErr
		}
	} else {
		balance, err := strconv.ParseInt(parameter.Balance, 10, 64)
		if err != nil {
			ccErr := GetError(1002)
			ccErr.AddMessage(err.Error())
			return nil, &ccErr
		}
		balanceInt = balance
	}

	appUser := AppUserStruct{
		Operation:         parameter.Operation,
		AppUserID:         parameter.AppUserID,
		TxTime:            parameter.TxTime,
		MessageHash:       parameter.MessageHash,
		ResearchProjectID: parameter.ResearchProjectID,
		Balance:           balanceInt,
		TransactionID:     stub.GetTxID(),
	}

	return appUser.Save(stub, parameter.Operation)
}

//GetAppUser :根据复合键进行前缀匹配查询,返回要查询的结构体的历史记录
//====================================================================
//args[0]: Operation 操作（复合主键）
//
//args[1]: ID 前缀匹配查询第一个参数
//
//args[...]: 允许传入更多的参数，可扩展。需按照存储时的前缀顺序传入更多参数
func GetAppUser(stub shim.ChaincodeStubInterface, args []string) ([]byte, *CcError) {
	var queryStructs []AppUserStruct
	resultsIterator, err := stub.GetStateByPartialCompositeKey(args[0], args)
	if err != nil {
		ccErr := GetError(1005)
		ccErr.AddMessage(err.Error())
		return nil, &ccErr
	}
	defer resultsIterator.Close()

	if resultsIterator.HasNext() == false {
		ccErr := GetError(1106)
		return nil, &ccErr
	}

	for resultsIterator.HasNext() {
		responseRange, err := resultsIterator.Next()
		if err != nil {
			ccErr := GetError(1106)
			ccErr.AddMessage(err.Error())
			return nil, &ccErr
		}

		if responseRange.Key != "" && responseRange.Value != nil {
			//fmt.Println("get responseRange.Key:======", responseRange.Key)
			//fmt.Println("get responseRange.Value:======", responseRange.Value)
			//拆分复合键
			_, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)
			if err != nil {
				ccErr := GetError(1006)
				ccErr.AddMessage(err.Error())
				return nil, &ccErr
			}

			var queryStruct AppUserStruct
			index, getIndexByOperationErr := GetIndexByOperation("AppUserStruct", args[0])
			if getIndexByOperationErr != nil {
				return nil, getIndexByOperationErr
			}

			queryStructMap := make(map[string]interface{})
			//fmt.Println("compositeKeyParts[index]===", compositeKeyParts[index])
			queryStructBytes := []byte(compositeKeyParts[index])
			errs := json.Unmarshal(queryStructBytes, &queryStructMap)
			if errs != nil {
				ccErr := GetError(1006)
				ccErr.AddMessage("entity get:" + errs.Error())
				return nil, &ccErr
			}

			errs = mapstructure.Decode(queryStructMap, &queryStruct)
			if errs != nil {
				ccErr := GetError(1013)
				ccErr.AddMessage("entity get: " + errs.Error())
				return nil, &ccErr
			}

			queryStructs = append(queryStructs, queryStruct)
		}
	}
	bytes, _ := json.Marshal(queryStructs)
	return bytes, nil
}

//QueryAppUserStruct :查询应用使用者结构体数据
//====================================================================
func QueryAppUserStruct(stub shim.ChaincodeStubInterface, parameter Parameter) ([]byte, *CcError) {
	var balanceInt int64
	balanceInt = 0
	if parameter.Balance != "" {
		balance, err := strconv.ParseInt(parameter.Balance, 10, 64)
		if err != nil {
			ccErr := GetError(1002)
			ccErr.AddMessage(err.Error())
			return nil, &ccErr
		}
		balanceInt = balance
	}

	appUser := AppUserStruct{
		Operation:         parameter.Operation,
		AppUserID:         parameter.AppUserID,
		TxTime:            parameter.TxTime,
		MessageHash:       parameter.MessageHash,
		Balance:           balanceInt,
		ResearchProjectID: parameter.ResearchProjectID,
		TransactionID:     parameter.TransactionID,
	}

	args, getAppUserStructPrefixErr := appUser.getAppUserStructPrefix(stub, parameter.Operation)
	if getAppUserStructPrefixErr != nil {
		return nil, getAppUserStructPrefixErr
	}
	var arg []string
	for _, value := range args[:len(args)-1] {
		if value != "" && value != "0" {
			arg = append(arg, value)
		} else {
			break
		}
	}
	
	bytes, errs := GetAppUser(stub, arg)
	if errs != nil || bytes == nil {
		return nil, errs
	}

	var appUserStructs, returnAppUserStruct,empty []AppUserStruct
	//如果开始和截止时间为空，则直接返回Get（）查询到的所有数据
	unmarshalErr := json.Unmarshal(bytes, &appUserStructs)
	if unmarshalErr != nil {
		ccErr := GetError(1006)
		ccErr.AddMessage("entity QueryAppUserStruct:" + unmarshalErr.Error())
		return nil, &ccErr
	}
	if parameter.StartTime != "" || parameter.EndTime != "" {
		for _, appUser := range appUserStructs {
			if parameter.StartTime <= appUser.TxTime && parameter.EndTime > appUser.TxTime {
				returnAppUserStruct = append(returnAppUserStruct, appUser)
			}
		}
	} else {
		returnAppUserStruct = appUserStructs
	}

	//分页
	var total, startCount, endCount, pages, resPerPage, resPage int64
	total, _ = strconv.ParseInt(strconv.Itoa(len(returnAppUserStruct)), 10, 64)

	//如果传入的分页中跳过的条数Page为空时，传回所有的数据
	if parameter.Page != "" {
		page, atoiErrs := strconv.ParseInt(parameter.Page, 10, 64)
		if atoiErrs != nil || page < 0 {
			ccErr := GetError(1002)
			logger.Warning(ccErr.GetJSON())
			return nil, &ccErr
		}
		resPage = page
		perPage, atoiErrs := strconv.ParseInt(parameter.PerPage, 10, 64)
		if atoiErrs != nil || perPage <= 0 {
			ccErr := GetError(1002)
			logger.Warning(ccErr.GetJSON())
			return nil, &ccErr
		}
		resPerPage = perPage

		startCount = (page - 1) * perPage //分页开始条数
		endCount = startCount + perPage   //分页截止条数

		if total > endCount {
			returnAppUserStruct = returnAppUserStruct[startCount:endCount]
		} else if total <= endCount && total > perPage {
			returnAppUserStruct = returnAppUserStruct[startCount:]
		} else {
			returnAppUserStruct = empty
		}
		if total%perPage != 0 {
			pages = total/perPage + 1
		} else {
			pages = total / perPage
		}
	}

	resAppUserStructs := ResponseAppUserStruct{
		Total:   total,
		Pages:   pages,
		PerPage: resPerPage,
		Page:    resPage,
		List:    returnAppUserStruct,
	}

	bytes, marshalErrs := json.Marshal(resAppUserStructs)
	if marshalErrs != nil {
		ccErr := GetError(1005)
		ccErr.AddMessage("entiyt QueryAppUserStruct :" + marshalErrs.Error())
		logger.Warning(ccErr.GetJSON())
		return nil, &ccErr
	}

	return bytes, nil
}
