package entity

import (
	"encoding/json"
	"github.com/goinggo/mapstructure"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strconv"
)

//Save :(it *RegulatorStruct)保存账户
//====================================================================
func (it *RegulatorStruct) Save(stub shim.ChaincodeStubInterface, compositeIndexName string) ([]byte, *CcError) {
	args, err := it.getRegulatorStructPrefix(stub, compositeIndexName)
	if err != nil {
		return nil, err
	}

	bytes, errs := Set(stub, args)
	if errs != nil {
		return nil, errs
	}
	return bytes, nil
}

//getRegulatorStructPrefix :(it *RegulatorStruct)查询该操作存储在链上的信息前缀
//====================================================================
func (it *RegulatorStruct) getRegulatorStructPrefix(stub shim.ChaincodeStubInterface, compositeIndexName string) ([]string, *CcError) {
	regulatorStruct := Struct2Map(*it)
	regulatorStructBytes, err := json.Marshal(regulatorStruct)
	if err != nil {
		ccErr := GetError(1006)
		ccErr.AddMessage("entity Save:" + err.Error())
		return nil, &ccErr
	}
	regulatorStructBytesStr := string(regulatorStructBytes)

	var args []string
	args = []string{compositeIndexName}
	switch compositeIndexName {
	case "newDatasetApproval", "modifyDatasetApproval":
		args = append(args, it.DataNodeID, it.DatasetID, it.Status, it.TxTime, it.MessageHash)
	case "researchProjectApproval", "useDatasetApproval":
		args = append(args, it.DataMiningID, it.ResearchProjectID, it.Status, it.TxTime, it.MessageHash)
	case "newAppApproval", "modifyAppApproval":
		args = append(args, it.AppProviderID, it.AppID, it.Status, it.TxTime, it.MessageHash)
	case "queryProjectResearchResults":
		args = append(args, it.RegulatorUserID, it.ResearchProjectID, it.TxTime, it.MessageHash)
	default:
		ccErr := GetError(1301)
		ccErr.AddMessage("entity getRegulatorStructPrefix")
		return nil, &ccErr
	}

	args = append(args, regulatorStructBytesStr)
	return args, nil
}

//RegulatorSave ：监管者数据写入
// ========================================================================
func RegulatorSave(stub shim.ChaincodeStubInterface, parameter Parameter) ([]byte, *CcError) {
	var datasets []DatasetStruct
	if parameter.Datasets != "" {
		errss := json.Unmarshal([]byte(parameter.Datasets), &datasets)
		if errss != nil {
			ccErr := GetError(1006)
			ccErr.AddMessage("entity QueryRegulatorStruct:" + errss.Error())
			return nil, &ccErr
		}
	}

	regulator := RegulatorStruct{
		Operation:         parameter.Operation,
		DataNodeID:        parameter.DataNodeID,
		NodeName:          parameter.NodeName,
		DatasetID:         parameter.DatasetID,
		TxTime:            parameter.TxTime,
		Status:            parameter.Status,
		MessageHash:       parameter.MessageHash,
		TransactionID:     stub.GetTxID(),
		DataMiningID:      parameter.DataMiningID,
		ResearchProjectID: parameter.ResearchProjectID,
		Datasets:          datasets,
		AppProviderID:     parameter.AppProviderID,
		AppID:             parameter.AppID,
		RegulatorUserID:   parameter.RegulatorUserID,
	}

	return regulator.Save(stub, parameter.Operation)
}

//GetRegulator :根据复合键进行前缀匹配查询,返回要查询的结构体的历史记录
//====================================================================
//args[0]: Operation 操作（复合主键）
//
//args[1]: ID 前缀匹配查询第一个参数
//
//args[...]: 允许传入更多的参数，可扩展。需按照存储时的前缀顺序传入更多参数
func GetRegulator(stub shim.ChaincodeStubInterface, args []string) ([]byte, *CcError) {
	var queryStructs []RegulatorStruct
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

			var queryStruct RegulatorStruct
			index, getIndexByOperationErr := GetIndexByOperation("RegulatorStruct", args[0])
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

//QueryRegulatorStruct :查询监管方结构体数据
//====================================================================
func QueryRegulatorStruct(stub shim.ChaincodeStubInterface, parameter Parameter) ([]byte, *CcError) {
	var datasetsEmpty []DatasetStruct
	regulator := RegulatorStruct{
		Operation:         parameter.Operation,
		DataNodeID:        parameter.DataNodeID,
		DatasetID:         parameter.DatasetID,
		TxTime:            parameter.TxTime,
		Status:            parameter.Status,
		MessageHash:       parameter.MessageHash,
		TransactionID:     parameter.TransactionID,
		DataMiningID:      parameter.DataMiningID,
		ResearchProjectID: parameter.ResearchProjectID,
		Datasets:          datasetsEmpty,
		AppProviderID:     parameter.AppProviderID,
		AppID:             parameter.AppID,
		RegulatorUserID:   parameter.RegulatorUserID,
	}

	args, err := regulator.getRegulatorStructPrefix(stub, parameter.Operation)
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

	regulators, errs := GetRegulator(stub, arg)
	if errs != nil || regulators == nil {
		return nil, errs
	}

	var RegulatorStructs, returnRegulatorStruct, empty []RegulatorStruct
	unmarshalErr := json.Unmarshal(regulators, &RegulatorStructs)
	if unmarshalErr != nil {
		ccErr := GetError(1006)
		ccErr.AddMessage("entity QueryRegulatorStruct:" + unmarshalErr.Error())
		return nil, &ccErr
	}

	if parameter.StartTime != "" || parameter.EndTime != "" {
		for _, Regulator := range RegulatorStructs {
			if parameter.StartTime <= Regulator.TxTime && parameter.EndTime > Regulator.TxTime {
				returnRegulatorStruct = append(returnRegulatorStruct, Regulator)
			}
		}
	} else {
		returnRegulatorStruct = RegulatorStructs
	}

	//分页
	var total, startCount, endCount, pages, resPerPage, resPage int64
	total, _ = strconv.ParseInt(strconv.Itoa(len(returnRegulatorStruct)), 10, 64)

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
			returnRegulatorStruct = returnRegulatorStruct[startCount:endCount]
		} else if total <= endCount && total > startCount {
			returnRegulatorStruct = returnRegulatorStruct[startCount:]
		} else {
			returnRegulatorStruct = empty
		}
		if total%perPage != 0 {
			pages = total/perPage + 1
		} else {
			pages = total / perPage
		}
	}

	resRegulatorStructs := ResponseRegulatorStruct{
		Total:   total,
		Pages:   pages,
		PerPage: resPerPage,
		Page:    resPage,
		List:    returnRegulatorStruct,
	}

	bytes, marshalErrs := json.Marshal(resRegulatorStructs)
	if marshalErrs != nil {
		ccErr := GetError(1005)
		ccErr.AddMessage("entiyt QueryRegulatorStruct :" + marshalErrs.Error())
		logger.Warning(ccErr.GetJSON())
		return nil, &ccErr
	}

	return bytes, nil
}
