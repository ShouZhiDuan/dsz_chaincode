package entity

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"reflect"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("entity")

//CatchErr :抓异常，并打印错误
func CatchErr(fun string, err interface{}) pb.Response {
	logger.Error(fun+"发生错误", err)

	ccErr := GetError(1202)
	logger.Warning(ccErr.GetJSON())
	return shim.Error(ccErr.GetJSON())
}

//CheckArgsLength  检查参数个数是否正确
func CheckArgsLength(args []string, num int) *CcError {
	if len(args) == num {
		return nil
	}
	ccErr := GetError(1001)
	argStr := "["
	for _, value := range args {
		argStr = argStr + " " + value + " "
	}
	argStr = argStr + "]"
	ccErr.AddMessage("Required " + strconv.Itoa(num) + " parameters , but get " + argStr)
	return &ccErr
}

//SortListByTime 按照时间排序
func SortListByTime(history []AppUserStruct) []AppUserStruct {
	num := len(history)
	for i := 0; i < num; i++ {
		for j := i + 1; j < num; j++ {
			if history[i].TxTime > history[j].TxTime {
				history[i], history[j] = history[j], history[i]
			}
		}
	}
	
	return history
}

//GetErrMessageByIndentity :根据用户身份获取错误提示信息
// ========================================================================
//--identity(string):	用户身份 包括：AppUser,AppProvider,DataUser,DataProvider,Regulator
func GetErrMessageByIndentity(identity, function string) string {
	var errMessage string
	switch identity {
	case "dataprovider":
		errMessage = "expecting [dataNodeRegistration||dataNodeModificationApply||addDatasetApply||openDatasets||modifyDatasets||datasetsOffline||queryResearchResults], but " + function
	case "appprovider":
		errMessage = "expecting [addAppApply||appOnline||modifyAppApply||queryAppUsage||appOffline], but " + function
	case "datauser":
		errMessage = "expecting [createResearchProjects||useDatasetApply||queryUserResearchResults||purchaseAuthorization||activateAuthorization], but " + function
	case "regulator":
		errMessage = "expecting [newDatasetApproval||modifyDatasetApproval||researchProjectApproval||useDatasetApproval||newAppApproval||modifyAppApproval||queryProjectResearchResults], but " + function
	case "appuser":
		errMessage = "expecting [newDatasetApproval||modifyDatasetApproval||researchProjectApproval||newAppApproval||modifyAppApproval||queryProjectResearchResults], but " + function
	default:
		errMessage = "the identity expecting [AppUser||AppProvider||DataUser||DataProvider||Regulator], but " + identity
	}
	return errMessage
}

//FilterParameter 处理参数
func FilterParameter(stub shim.ChaincodeStubInterface, args []string, invokeOrQuery string) (Parameter, *CcError) {
	var empty Parameter
	ccErr := CheckArgsLength(args, 1)
	if ccErr != nil {
		logger.Warning(ccErr.GetJSON())
		return empty, ccErr
	}

	var parameter Parameter
	errss1 := json.Unmarshal([]byte(args[0]), &parameter)
	if errss1 != nil {
		ccErr := GetError(1012)
		logger.Warning(ccErr.GetJSON())
		return empty, &ccErr
	}                                 

	if invokeOrQuery == "invoke" && parameter.TxTime == "" {
		parameter.TxTime = CheckTimeStamp(stub)
	}

	return parameter, nil
}

//CheckArgsNil 检查参数是否为空字符串
func CheckArgsNil(args []string) *CcError {
	for i, x := range args {
		if x == "" {
			ccErr := GetError(1015)
			ccErr.AddMessage("the number of " + strconv.Itoa(i) + " parameter is empty")
			return &ccErr
		}
	}
	return nil
}

//CheckTimeStamp 返回当前时间:20190614091705
func CheckTimeStamp(stub shim.ChaincodeStubInterface) string {
	l, _ := time.LoadLocation("Asia/Shanghai")

	timeStamp, _ := stub.GetTxTimestamp()
	time0 := time.Unix(timeStamp.Seconds, int64(timeStamp.Nanos))

	times := time0.In(l).Format("20060102150405")
	// logger.Infoln("timeStamp转换成字符串的上海时区时间", time4) //20190703111541
	return times
}

//Struct2Map :结构体转为map
func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

//InvokeOtherChainCode :
// ========================================================================
//调用其他链上chaincode
// ========================================================================
//--stub(shim.ChaincodeStubInterface):链码接口
//
//--args[0]: Channel_name
//
//--args[1]: Chaincode_name 实例名
//
//--args[2]: func 方法名
//
//--args[3]: args[...] 参数
func InvokeOtherChainCode(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//logger.Infoln("args:", args)
	var params []string
	for j, param := range args {
		if j > 1 {
			params = append(params, param)
		}
	}

	//string数组转化为byte数组
	queryArgs := make([][]byte, len(params))
	for i, arg := range params {
		queryArgs[i] = []byte(arg)
	}

	//返回data数据
	response := stub.InvokeChaincode(args[1], queryArgs, args[0])
	if response.Status != shim.OK {
		errStr := GetError(1302)
		return shim.Error(errStr.GetJSON())
	}

	return shim.Success(response.Payload)
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////

//TimeFormatStringToTime 时间格式转换，string类型转换成time.Time类型
func TimeFormatStringToTime(str string) (time.Time, *CcError) {
	timeLayout := "20060102150405"
	loc, _ := time.LoadLocation("Local")
	theTime, err := time.ParseInLocation(timeLayout, str, loc)
	if err != nil {
		ccErr := GetError(1017)
		return theTime, &ccErr
	}
	return theTime, nil
}

//GetUserName :获取调用链码的用户名
// ========================================================================
//--stub(shim.ChaincodeStubInterface):	链码接口
func GetUserName(stub shim.ChaincodeStubInterface) string {
	creatorByte, _ := stub.GetCreator()
	certStart := bytes.IndexAny(creatorByte, "-----BEGIN")
	if certStart == -1 {
		logger.Info("No certificate found")
		return ""
	}
	certText := creatorByte[certStart:]
	bl, _ := pem.Decode(certText)
	if bl == nil {
		logger.Info("Could not decode the PEM structure")
		return ""
	}

	cert, err := x509.ParseCertificate(bl.Bytes)
	if err != nil {
		logger.Info("ParseCertificate failed")
		return ""
	}
	uname := cert.Subject.CommonName
	return uname
}
