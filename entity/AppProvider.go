package entity

import (
	"encoding/json"
	"github.com/goinggo/mapstructure"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strconv"
)

//Save :(it *AppProviderStruct)保存账户
//====================================================================
//compositeIndexName 复合键
func (it *AppProviderStruct) Save(stub shim.ChaincodeStubInterface, compositeIndexName string) ([]byte, *CcError) {
	args, err := it.getAppProviderStructPrefix(stub, compositeIndexName)
	if err != nil {
		return nil, err
	}

	bytes, errs := Set(stub, args)
	if errs != nil {
		return nil, errs
	}
	return bytes, nil
}

//getAppProviderStructPrefix :(it *AppProviderStruct)查询该操作存储在链上的信息前缀
//====================================================================
func (it *AppProviderStruct) getAppProviderStructPrefix(stub shim.ChaincodeStubInterface, compositeIndexName string) ([]string, *CcError) {
	appProviderStruct := Struct2Map(*it)
	appProviderStructBytes, err := json.Marshal(appProviderStruct)
	if err != nil {
		ccErr := GetError(1006)
		ccErr.AddMessage("entity Save:" + err.Error())
		return nil, &ccErr
	}
	appProviderStructBytesStr := string(appProviderStructBytes)

	var args []string
	args = []string{compositeIndexName}
	switch compositeIndexName {
	case "addAppApply", "appOnline", "modifyAppApply", "appOffline":
		args = append(args, it.AppProviderID, it.AppID, it.Status, it.TxTime, it.MessageHash)
	case "queryAppUsage":
		args = append(args, it.AppProviderID, it.AppID, it.TxTime, it.MessageHash)
	default:
		ccErr := GetError(1301)
		ccErr.AddMessage("entity getAppProviderStructPrefix")
		return nil, &ccErr
	}

	args = append(args, appProviderStructBytesStr)
	return args, nil
}

//AppProviderSave ：应用提供方数据写入
// ========================================================================
func AppProviderSave(stub shim.ChaincodeStubInterface, parameter Parameter) ([]byte, *CcError) {
	appProvider := AppProviderStruct{
		Operation:     parameter.Operation,
		AppProviderID: parameter.AppProviderID,
		AppID:         parameter.AppID,
		TxTime:        parameter.TxTime,
		Status:        parameter.Status,
		MessageHash:   parameter.MessageHash,
		TransactionID: stub.GetTxID(),
	}

	return appProvider.Save(stub, parameter.Operation)
}

//GetAppProvider :根据复合键进行前缀匹配查询,返回要查询的结构体的历史记录
//====================================================================
//args[0]: Operation 操作（复合主键）
//
//args[1]: ID 前缀匹配查询第一个参数
//
//args[...]: 允许传入更多的参数，可扩展。需按照存储时的前缀顺序传入更多参数
func GetAppProvider(stub shim.ChaincodeStubInterface, args []string) ([]byte, *CcError) {
	var queryStructs []AppProviderStruct
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

			var queryStruct AppProviderStruct
			index, getIndexByOperationErr := GetIndexByOperation("AppProviderStruct", args[0])
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

			errs = mapstructure.Decode(queryStructMap, &queryStruct) //map[string][]interface{}
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

//QueryAppProviderStruct :查询应用提供方结构体数据
//====================================================================
func QueryAppProviderStruct(stub shim.ChaincodeStubInterface, parameter Parameter) ([]byte, *CcError) {
	appProvider := AppProviderStruct{
		Operation:     parameter.Operation,
		AppProviderID: parameter.AppProviderID,
		AppID:         parameter.AppID,
		TxTime:        parameter.TxTime,
		Status:        parameter.Status,
		MessageHash:   parameter.MessageHash,
		TransactionID: parameter.TransactionID,
	}

	args, err := appProvider.getAppProviderStructPrefix(stub, parameter.Operation)
	if err != nil {
		return nil, err
	}
	var arg []string
	for _, value := range args[:len(args)-1] {
		if value != "" {
			arg = append(arg, value)
		} else {
			break
		}
	}

	bytes, errs := GetAppProvider(stub, arg)
	if errs != nil || bytes == nil {
		return nil, errs
	}

	var appProviderStructs, returnAppProviderStruct,empty []AppProviderStruct
	//如果开始和截止时间为空，则直接返回Get（）查询到的所有数据
	unmarshalErr := json.Unmarshal(bytes, &appProviderStructs)
	if unmarshalErr != nil {
		ccErr := GetError(1006)
		ccErr.AddMessage("entity QueryAppProviderStruct:" + unmarshalErr.Error())
		return nil, &ccErr
	}
	if parameter.StartTime != "" || parameter.EndTime != "" {
		for _, appProvider := range appProviderStructs {
			if parameter.StartTime <= appProvider.TxTime && parameter.EndTime > appProvider.TxTime {
				returnAppProviderStruct = append(returnAppProviderStruct, appProvider)
			}
		}
	} else {
		returnAppProviderStruct = appProviderStructs
	}

	//分页
	var total, startCount, endCount, pages, resPerPage, resPage int64
	total, _ = strconv.ParseInt(strconv.Itoa(len(returnAppProviderStruct)), 10, 64)

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
			returnAppProviderStruct = returnAppProviderStruct[startCount:endCount]
		} else if total <= endCount && total > perPage {
			returnAppProviderStruct = returnAppProviderStruct[startCount:]
		} else {
			returnAppProviderStruct = empty
		}
		if total%perPage != 0 {
			pages = total/perPage + 1
		} else {
			pages = total / perPage
		}
	}

	resAppProviderStructs := ResponseAppProviderStruct{
		Total:   total,
		Pages:   pages,
		PerPage: resPerPage,
		Page:    resPage,
		List:    returnAppProviderStruct,
	}
	
	bytes, marshalErrs := json.Marshal(resAppProviderStructs)
	if marshalErrs != nil {
		ccErr := GetError(1005)
		ccErr.AddMessage("entiyt QueryAppProviderStruct :" + marshalErrs.Error())
		logger.Warning(ccErr.GetJSON())
		return nil, &ccErr
	}

	return bytes, nil
}
