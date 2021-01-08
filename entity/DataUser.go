package entity

import (
	"encoding/json"
	"github.com/goinggo/mapstructure"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strconv"
)

//Save :(it *DataUserStruct)保存账户
//====================================================================
//compositeIndexName 复合键
func (it *DataUserStruct) Save(stub shim.ChaincodeStubInterface, compositeIndexName string) ([]byte, *CcError) {
	args, err := it.getDataUserStructPrefix(stub, compositeIndexName)
	if err != nil {
		return nil, err
	}

	bytes, errs := Set(stub, args)
	if errs != nil {
		return nil, errs
	}
	return bytes, nil
}

//getDataUserStructPrefix :(it *DataUserStruct)查询该操作存储在链上的信息前缀
//====================================================================
//compositeIndexName 复合键
func (it *DataUserStruct) getDataUserStructPrefix(stub shim.ChaincodeStubInterface, compositeIndexName string) ([]string, *CcError) {
	dataUserStruct := Struct2Map(*it)
	dataUserStructBytes, err := json.Marshal(dataUserStruct)
	if err != nil {
		ccErr := GetError(1006)
		ccErr.AddMessage("entity Save:" + err.Error())
		return nil, &ccErr
	}
	dataUserStructBytesStr := string(dataUserStructBytes)

	var args []string
	args = []string{compositeIndexName}
	switch compositeIndexName {
	case "createResearchProjects", "useDatasetApply":
		args = append(args, it.DataMiningID, it.ResearchProjectID, it.Status, it.TxTime, it.MessageHash)
	case "queryUserResearchResults":
		args = append(args, it.DataMiningID, it.ResearchProjectID, it.TxTime, it.MessageHash)
	case "purchaseAuthorization", "activateAuthorization":
		args = append(args, it.DataMiningID, it.CertificateID, it.TxTime, it.MessageHash)
	default:
		ccErr := GetError(1301)
		ccErr.AddMessage("entity getDataUserStructPrefix")
		return nil, &ccErr
	}

	args = append(args, dataUserStructBytesStr)
	return args, nil
}

//DataUserSave ：数据使用者数据写入
// ========================================================================
func DataUserSave(stub shim.ChaincodeStubInterface, parameter Parameter) ([]byte, *CcError) {
	var datasets []DatasetStruct
	if parameter.Datasets != "" {
		errss := json.Unmarshal([]byte(parameter.Datasets), &datasets)
		if errss != nil {
			ccErr := GetError(1006)
			ccErr.AddMessage("entity QueryDataUserStruct:" + errss.Error())
			return nil, &ccErr
		}
	}
	
	dataUser := DataUserStruct{
		Operation:         parameter.Operation,
		DataMiningID:      parameter.DataMiningID,
		ResearchProjectID: parameter.ResearchProjectID,
		TxTime:            parameter.TxTime,
		MessageHash:       parameter.MessageHash,
		Status:            parameter.Status,
		Datasets:          datasets, 
		CertificateID:     parameter.CertificateID,
		TransactionID:     stub.GetTxID(),
	}

	return dataUser.Save(stub, parameter.Operation)
}

//GetDataUser :根据复合键进行前缀匹配查询,返回要查询的结构体的历史记录
//====================================================================
//args[0]: Operation 操作（复合主键）
//
//args[1]: ID 前缀匹配查询第一个参数
//
//args[...]: 允许传入更多的参数，可扩展。需按照存储时的前缀顺序传入更多参数
func GetDataUser(stub shim.ChaincodeStubInterface, args []string) ([]byte, *CcError) {
	var queryStructs []DataUserStruct
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

			var queryStruct DataUserStruct
			index, getIndexByOperationErr := GetIndexByOperation("DataUserStruct", args[0])
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

//QueryDataUserStruct :查询数据使用者结构体数据
//====================================================================
func QueryDataUserStruct(stub shim.ChaincodeStubInterface, parameter Parameter) ([]byte, *CcError) {
	var datasetsEmpty []DatasetStruct
	dataUser := DataUserStruct{
		Operation:         parameter.Operation,
		DataMiningID:      parameter.DataMiningID,
		ResearchProjectID: parameter.ResearchProjectID,
		TxTime:            parameter.TxTime,
		MessageHash:       parameter.MessageHash,
		Status:            parameter.Status,
		Datasets:          datasetsEmpty,
		CertificateID:     parameter.CertificateID,
		TransactionID:     parameter.TransactionID,
	}

	args, err := dataUser.getDataUserStructPrefix(stub, parameter.Operation)
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

	dataUsers, errs := GetDataUser(stub, arg)
	if errs != nil || dataUsers == nil {
		return nil, errs
	}

	var dataUserStructs, returnDataUserStruct,empty []DataUserStruct
	unmarshalErr := json.Unmarshal(dataUsers, &dataUserStructs)
	if unmarshalErr != nil {
		ccErr := GetError(1006)
		ccErr.AddMessage("entity QueryDataUserStruct:" + unmarshalErr.Error())
		return nil, &ccErr
	}

	if parameter.StartTime != "" || parameter.EndTime != "" {
		for _, dataUser := range dataUserStructs {
			if parameter.StartTime <= dataUser.TxTime && parameter.EndTime > dataUser.TxTime {
				returnDataUserStruct = append(returnDataUserStruct, dataUser)
			}
		}
	} else {
		returnDataUserStruct = dataUserStructs
	}

	//分页
	var total, startCount, endCount, pages, resPerPage, resPage int64
	total, _ = strconv.ParseInt(strconv.Itoa(len(returnDataUserStruct)), 10, 64)

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
			returnDataUserStruct = returnDataUserStruct[startCount:endCount]
		} else if total <= endCount && total > startCount {
			returnDataUserStruct = returnDataUserStruct[startCount:]
		} else {
			returnDataUserStruct = empty
		}
		if total%perPage != 0 {
			pages = total/perPage + 1
		} else {
			pages = total / perPage
		}
	}

	resDataUserStructs := ResponseDataUserStruct{
		Total:   total,
		Pages:   pages,
		PerPage: resPerPage,
		Page:    resPage,
		List:    returnDataUserStruct,
	}

	bytes, marshalErrs := json.Marshal(resDataUserStructs)
	if marshalErrs != nil {
		ccErr := GetError(1005)
		ccErr.AddMessage("entiyt QueryDataUserStruct :" + marshalErrs.Error())
		logger.Warning(ccErr.GetJSON())
		return nil, &ccErr
	}

	return bytes, nil
}
