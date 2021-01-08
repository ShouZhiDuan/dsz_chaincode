package entity

import (
	"encoding/json"
	"github.com/goinggo/mapstructure"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strconv"
)

//Set :创建复合键
//====================================================================
//args[0]: compositeIndexName 复合键
//
//args[1]: key 前缀匹配查询第一个参数
//
//args[2]:  结构体对象转化为的字符串数组（这里不做判断，只是存储）
//
//args[...]: 允许存储更多的参数，可扩展
func Set(stub shim.ChaincodeStubInterface, args []string) ([]byte, *CcError) {
	compositeIndexName := args[0]
	compositeKey, compositeErr := stub.CreateCompositeKey(compositeIndexName, args)
	if compositeErr != nil {
		ccErr := GetError(1009)
		ccErr.AddMessage(compositeErr.Error())
		return nil, &ccErr
	}

	compositePutErr := stub.PutState(compositeKey, []byte{0x00})
	if compositePutErr != nil {
		ccErr := GetError(1010)
		ccErr.AddMessage(compositePutErr.Error())
		return nil, &ccErr
	}

	bytes, _ := json.Marshal(args)
	return bytes, nil
}

//deleteData :删除数据
//====================================================================
// key: 区块链里的key
func DeleteData(stub shim.ChaincodeStubInterface, key string) *CcError {
	err := stub.DelState(key)
	if err != nil {
		ccErr := GetError(1102)
		ccErr.AddMessage("entity deleteData ：" + err.Error())
		return &ccErr
	}
	return nil
}

//Save :(it *DataProviderStruct)保存账户
//====================================================================
//compositeIndexName 复合键
func (it *DataProviderStruct) Save(stub shim.ChaincodeStubInterface, compositeIndexName string) ([]byte, *CcError) {
	args, err := it.getDataProviderStructPrefix(stub, compositeIndexName)
	if err != nil {
		return nil, err
	}

	bytes, errs := Set(stub, args)
	if errs != nil {
		return nil, errs
	}
	return bytes, nil
}

//DataProviderSave ：数据提供者数据写入
// ========================================================================
func DataProviderSave(stub shim.ChaincodeStubInterface, parameter Parameter) ([]byte, *CcError) {
	dataProvider := DataProviderStruct{
		Operation:     parameter.Operation,
		DataNodeID:    parameter.DataNodeID,
		TxTime:        parameter.TxTime,
		NodeName:      parameter.NodeName,
		MessageHash:   parameter.MessageHash,
		DatasetID:     parameter.DatasetID,
		Status:        parameter.Status,
		TransactionID: stub.GetTxID(),
	}

	return dataProvider.Save(stub, parameter.Operation)
}

//getDataProviderStructPrefix :(it *DataProviderStruct)查询该操作存储在链上的信息前缀
//====================================================================
//compositeIndexName 复合键
func (it *DataProviderStruct) getDataProviderStructPrefix(stub shim.ChaincodeStubInterface, compositeIndexName string) ([]string, *CcError) {
	dataProviderStruct := Struct2Map(*it)
	dataProviderStructBytes, err := json.Marshal(dataProviderStruct)
	if err != nil {
		ccErr := GetError(1006)
		ccErr.AddMessage("entity Save:" + err.Error())
		return nil, &ccErr
	}
	dataProviderStructBytesStr := string(dataProviderStructBytes)

	var args []string
	args = []string{compositeIndexName}
	switch compositeIndexName {
	//      数据节点注册          数据节点修改申请
	case "dataNodeRegistration", "dataNodeModificationApply":
		args = append(args, it.DataNodeID, it.NodeName, it.TxTime, it.MessageHash)
	case "addDatasetApply", "openDatasets", "modifyDatasets", "datasetsOffline":
		args = append(args, it.DataNodeID, it.NodeName, it.DatasetID, it.Status, it.TxTime, it.MessageHash)
	case "queryNodeDataResearchResults":
		args = append(args, it.DataNodeID, it.NodeName, it.DatasetID, it.TxTime, it.MessageHash)
	default:
		ccErr := GetError(1301)
		ccErr.AddMessage("entity getDataProviderStructPrefix")
		return nil, &ccErr
	}

	args = append(args, dataProviderStructBytesStr)
	return args, nil
}

//GetIndexByOperation :查询该操作存储在链上的信息索引
//====================================================================
func GetIndexByOperation(structName string, operation string) (int, *CcError) {
	var index int
	switch structName {
	case "DataProviderStruct":
		switch operation {
		case "dataNodeRegistration", "dataNodeModificationApply":
			index = 5
		case "addDatasetApply", "openDatasets", "modifyDatasets", "datasetsOffline":
			index = 7
		case "queryNodeDataResearchResults":
			index = 6
		default:
			ccErr := GetError(1301)
			ccErr.AddMessage("entity getDataProviderStructPrefix DataProviderStruct")
			return 0, &ccErr
		}
	case "AppProviderStruct":
		switch operation {
		case "addAppApply", "appOnline", "modifyAppApply", "appOffline":
			index = 6
		case "queryAppUsage":
			index = 5
		default:
			ccErr := GetError(1301)
			ccErr.AddMessage("entity getDataProviderStructPrefix AppProviderStruct")
			return 0, &ccErr
		}
	case "DataUserStruct":
		switch operation {
		case "createResearchProjects", "useDatasetApply":
			index = 6
		case "queryUserResearchResults", "purchaseAuthorization", "activateAuthorization":
			index = 5
		default:
			ccErr := GetError(1301)
			ccErr.AddMessage("entity getDataProviderStructPrefix DataUserStruct")
			return 0, &ccErr
		}
	case "RegulatorStruct":
		switch operation {
		case "newDatasetApproval", "modifyDatasetApproval", "researchProjectApproval", "useDatasetApproval", "newAppApproval", "modifyAppApproval":
			index = 6
		case "queryProjectResearchResults":
			index = 5
		default:
			ccErr := GetError(1301)
			ccErr.AddMessage("entity getDataProviderStructPrefix RegulatorStruct")
			return 0, &ccErr
		}
	case "AppUserStruct":
		switch operation {
		case "authorizationDataUpBlockChain", "queryAppUserDataResearchResults":
			index = 5
		case "authorizationAndParticipationInResearch":
			index = 6
		case "queryDataUsage":
			index = 4
		default:
			ccErr := GetError(1301)
			ccErr.AddMessage("entity getDataProviderStructPrefix AppUserStruct")
			return 0, &ccErr
		}
	default:
		ccErr := GetError(1301)
		ccErr.AddMessage("entity getDataProviderStructPrefix")
		return 0, &ccErr
	}

	return index, nil
}

//GetDataProvider :根据复合键进行前缀匹配查询,返回要查询的结构体的历史记录
//====================================================================
//args[0]: Operation 操作（复合主键）
//
//args[1]: ID 前缀匹配查询第一个参数
//
//args[...]: 允许传入更多的参数，可扩展。需按照存储时的前缀顺序传入更多参数
func GetDataProvider(stub shim.ChaincodeStubInterface, args []string) ([]byte, *CcError) {
	var queryStructs []DataProviderStruct
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
				ccErr := GetError(1012)
				ccErr.AddMessage(err.Error())
				return nil, &ccErr
			}

			var queryStruct DataProviderStruct
			index, getIndexByOperationErr := GetIndexByOperation("DataProviderStruct", args[0])
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

//QueryDataProviderStruct :查询数据提供方结构体数据
//====================================================================
func QueryDataProviderStruct(stub shim.ChaincodeStubInterface, parameter Parameter) ([]byte, *CcError) {
	dataProvider := DataProviderStruct{
		Operation:     parameter.Operation,
		DataNodeID:    parameter.DataNodeID,
		TxTime:        parameter.TxTime,
		NodeName:      parameter.NodeName,
		MessageHash:   parameter.MessageHash,
		DatasetID:     parameter.DatasetID,
		Status:        parameter.Status,
		TransactionID: parameter.TransactionID,
	}

	args, err := dataProvider.getDataProviderStructPrefix(stub, parameter.Operation)
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

	dataProviders, errs := GetDataProvider(stub, arg)
	if errs != nil || dataProviders == nil {
		return nil, errs
	}

	var dataProviderStructs, returnDataProviderStruct, empty []DataProviderStruct
	unmarshalErr := json.Unmarshal(dataProviders, &dataProviderStructs)
	if unmarshalErr != nil {
		ccErr := GetError(1006)
		ccErr.AddMessage("entity QueryDataProviderStruct:" + unmarshalErr.Error())
		return nil, &ccErr
	}

	if parameter.StartTime != "" || parameter.EndTime != "" {
		for _, dataProvider := range dataProviderStructs {
			if parameter.StartTime <= dataProvider.TxTime && parameter.EndTime > dataProvider.TxTime {
				returnDataProviderStruct = append(returnDataProviderStruct, dataProvider)
			}
		}
	} else {
		returnDataProviderStruct = dataProviderStructs
	}

	//分页
	var total, startCount, endCount, pages, resPerPage, resPage int64
	total, _ = strconv.ParseInt(strconv.Itoa(len(returnDataProviderStruct)), 10, 64)

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
			returnDataProviderStruct = returnDataProviderStruct[startCount:endCount]
		} else if total <= endCount && total > startCount {
			returnDataProviderStruct = returnDataProviderStruct[startCount:]
		} else {
			returnDataProviderStruct = empty
		}
		if total%perPage != 0 {
			pages = total/perPage + 1
		} else {
			pages = total / perPage
		}
	}

	resDataProviderStructs := ResponseDataProviderStruct{
		Total:   total,
		Pages:   pages,
		PerPage: resPerPage,
		Page:    resPage,
		List:    returnDataProviderStruct,
	}

	bytes, marshalErrs := json.Marshal(resDataProviderStructs)
	if marshalErrs != nil {
		ccErr := GetError(1005)
		ccErr.AddMessage("entiyt QueryDataProviderStruct :" + marshalErrs.Error())
		logger.Warning(ccErr.GetJSON())
		return nil, &ccErr
	}

	return bytes, nil
}
