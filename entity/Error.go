package entity

import "encoding/json"

//CcError 错误的结构体
// ========================================================================
type CcError struct {
	ErrNO   int    //错误序号
	Message string //错误内容
}

//Errors 所有错误序号和对应错误内容
// ========================================================================
var Errors = map[int]string{
	//参数类错误
	1000: "Parametric class error",                //参数类错误
	1001: "Incorrect number of parameters",        //参数个数不正确
	1002: "Incorrect parameter type",              //参数类型不正确
	1003: "Iterative data incorrect",              //迭代数据不正确
	1004: "Data conversion failed",                //数据转换失败
	1005: "Incorrect marshalling of status data",  //状态数据编组不正确
	1006: "Incorrect unmarshalling of state data", //状态数据解组不正确
	1007: "The afferent parameter is empty",       //参数为空
	1008: "Time Format error",                     //时间格式转换错误
	1009: "Failed to create compositeKey",         //创建复合键失败
	1010: "Failed to store compositeKey",          //存储复合键失败
	1011: "Failed to query through compositeKey",  //通过复合键查询失败
	1012: "Failed to split compositeKey",          //拆分复合键失败
	1013: "Map type to struct type failed",        //map类型转struct类型失败

	//世界状态内错误
	1100: "Errors in the World State ",                          //世界状态内错误
	1101: "Failure of state data storage from the World State",  //状态数据存储失败
	1102: "No status data was found from the World State",       //找不到状态数据
	1103: "The rich text query failed from the World State",     //富文本查询失败
	1104: "Failed to delete data from leveldb",                  //从leveldb数据库中删除数据失败
	1105: "The Data already exists and cannot be created again", //数据已经存在，不能再次被创建
	1106: "No historical data was found",                        //找不到历史数据

	//链上错误
	1200: "Chain error",          //链上错误
	1201: "Block data not found", //找不到块数据
	1202: "Errors in goroutine",  //运行时错误

	//chaincode错误
	1300: "Chaincode error",                        //chaincode错误
	1301: "Can not find the method from Chaincode", //找不到方法
	1302: "Invoke other Chaincode error",           //调用其他chaincode错误

	//业务错误

}

//GetError 返回错误的对象
// ========================================================================
func GetError(errno int) CcError {
	return CcError{ErrNO: errno, Message: Errors[errno]}
}

// AddMessage 添加错误消息
// ========================================================================
//--message（string）: 附加到错误文字的后面
func (it *CcError) AddMessage(message string) {
	it.Message = it.Message + ":" + message
}

//Get 返回错误的对象
// ========================================================================
//--errno（int）: 错误序号
func (it *CcError) Get(errno int) CcError {
	return CcError{ErrNO: errno, Message: Errors[errno]}
}

//GetJSON 返回错误的json字符串
// ========================================================================
func (it *CcError) GetJSON() string {
	bytes, _ := json.Marshal(it)
	return string(bytes)
}
