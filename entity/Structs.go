package entity

const (
	//状态：已申请
	StatusApplied = "Applied"
	//状态：审批通过
	StatusApproved = "Approved"
	//状态：审批未通过
	StatusApprovedFailed = "ApprovalFailed"
	//状态：应用上架
	StatusOnline = "Online"
	//状态：应用下架
	StatusOffline = "Offline"
)

//Parameter 传入参数结构体  (传入json格式的字符串或者[]string数组)
type Parameter struct {
	ChannelName   string `json:"channelName"`   //通道名称
	ChaincodeName string `json:"chaincodeName"` //链码名称
	FunctionName  string `json:"functionName"`  //方法名称
	StartTime     string `json:"startTime"`     //开始时间
	EndTime       string `json:"endTime"`       //截止时间
	Page          string `json:"page"`          //当前页数
	PerPage       string `json:"perPage"`       //每页几条
	TransactionID string `json:"transactionID"` //该交易在链上的TransactionID（在SDK端，可通过该ID查询指定的交易）
	Operation     string `json:"operation"`     //操作

	DataNodeID  string `json:"dataNodeID"`  //数据节点ID
	TxTime      string `json:"txTime"`      //交易创建时间
	NodeName    string `json:"nodeName"`    //节点名称
	MessageHash string `json:"messageHash"` //信息Hash
	DatasetID   string `json:"datasetID"`   //数据集ID
	Status      string `json:"status"`      //数据集状态

	AppProviderID string `json:"appProviderStruct"` //应用提供方ID
	AppID         string `json:"appID"`             //应用ID

	DataMiningID      string `json:"dataMiningID"`      //数据挖掘者ID
	ResearchProjectID string `json:"researchProjectID"` //研究项目ID
	Datasets          string `json:"datasets"`          //数据集  （需要调接口时调试如何传值）
	CertificateID     string `json:"certificateID"`     //证书ID
	DatasetHash       string `json:"datasetHash"`       //数据集Hash

	RegulatorUserID string `json:"regulatorUserID"` //监管者用户ID

	AppUserID string `json:"appUserID"` //应用用户ID
	Balance   string `json:"Balance"`   //积分
}

// ResponseDataProviderStruct  查询返回数据提供方列表
type ResponseDataProviderStruct struct {
	Total   int64                `json:"total"`   //总条数
	Pages   int64                `json:"pages"`   //总页数
	PerPage int64                `json:"perPage"` //每页条数
	Page    int64               `json:"page"`    //当前页数
	List    []DataProviderStruct `json:"list"`    //数据提供方结构体列表
}

//DataProviderStruct 数据提供方结构体
type DataProviderStruct struct {
	Operation     string `json:"operation"`     //操作
	DataNodeID    string `json:"dataNodeID"`    //数据节点ID
	TxTime        string `json:"txTime"`        //交易创建时间
	NodeName      string `json:"nodeName"`      //节点名称
	MessageHash   string `json:"messageHash"`   //信息Hash
	DatasetID     string `json:"datasetID"`     //数据集ID
	Status        string `json:"status"`        //数据集状态
	TransactionID string `json:"transactionID"` //该交易在链上的TransactionID
}

// ResponseAppProviderStruct  查询返回应用提供方列表
type ResponseAppProviderStruct struct {
	Total   int64               `json:"total"`   //总条数
	Pages   int64               `json:"pages"`   //总页数
	PerPage int64               `json:"perPage"` //每页条数
	Page    int64              `json:"page"`    //当前页数
	List    []AppProviderStruct `json:"list"`    //应用提供方结构体列表
}

//AppProviderStruct 应用提供方结构体
type AppProviderStruct struct {
	Operation     string `json:"operation"`         //操作
	AppProviderID string `json:"appProviderStruct"` //应用提供方ID
	AppID         string `json:"appID"`             //应用ID
	TxTime        string `json:"txTime"`            //交易创建时间
	Status        string `json:"status"`            //应用状态
	MessageHash   string `json:"messageHash"`       //信息hash
	TransactionID string `json:"transactionID"`     //该交易在链上的TransactionID
}

// ResponseDataUserStruct  查询返回数据使用方列表
type ResponseDataUserStruct struct {
	Total   int64            `json:"total"`   //总条数
	Pages   int64            `json:"pages"`   //总页数
	PerPage int64            `json:"perPage"` //每页条数
	Page    int64           `json:"page"`    //当前页数
	List    []DataUserStruct `json:"list"`    //数据使用方结构体列表
}

//DataUserStruct 数据使用方结构体
type DataUserStruct struct {
	Operation         string          `json:"operation"`         //操作
	DataMiningID      string          `json:"dataMiningID"`      //数据挖掘者ID
	ResearchProjectID string          `json:"researchProjectID"` //研究项目ID
	TxTime            string          `json:"txTime"`            //交易创建时间
	MessageHash       string          `json:"messageHash"`       //信息hash
	Datasets          []DatasetStruct `json:"datasets"`          //数据集
	Status            string          `json:"status"`            //状态
	CertificateID     string          `json:"certificateID"`     //证书ID
	TransactionID     string          `json:"transactionID"`     //该交易在链上的TransactionID
}

//DatasetStruct 数据集结构体
type DatasetStruct struct {
	DataNodeID  string `json:"dataNodeID"`  //数据节点ID
	DatasetID   string `json:"datasetID"`   //数据集ID
	DatasetHash string `json:"datasetHash"` //数据集Hash
}

// ResponseRegulatorStruct  查询返回监管方列表
type ResponseRegulatorStruct struct {
	Total   int64             `json:"total"`   //总条数
	Pages   int64             `json:"pages"`   //总页数
	PerPage int64             `json:"perPage"` //每页条数
	Page    int64            `json:"page"`    //当前页数
	List    []RegulatorStruct `json:"list"`    //监管方结构体列表
}

//RegulatorStruct 监管方结构体
type RegulatorStruct struct {
	Operation     string `json:"operation"`     //操作
	DataNodeID    string `json:"dataNodeID"`    //数据节点ID
	NodeName      string `json:"nodeName"`      //节点名称
	DatasetID     string `json:"datasetId"`     //数据集ID
	TxTime        string `json:"txTime"`        //审批时间
	Status        string `json:"status"`        //状态
	MessageHash   string `json:"messageHash"`   //信息Hash
	TransactionID string `json:"transactionID"` //该交易在链上的TransactionID

	DataMiningID      string          `json:"dataMiningID"`      //数据挖掘者ID
	ResearchProjectID string          `json:"researchProjectID"` //研究项目ID
	Datasets          []DatasetStruct `json:"datasets"`          //数据集

	AppProviderID string `json:"appProviderID"` //应用提供方ID
	AppID         string `json:"appID"`         //应用ID

	RegulatorUserID string `json:"regulatorUserID"` //监管者用户ID
}

// ResponseAppUserStruct  查询返回应用用户列表
type ResponseAppUserStruct struct {
	Total   int64           `json:"total"`   //总条数
	Pages   int64           `json:"pages"`   //总页数
	PerPage int64           `json:"perPage"` //每页条数
	Page    int64          `json:"page"`    //当前页数
	List    []AppUserStruct `json:"list"`    //应用用户结构体列表
}

//AppUserStruct 应用用户结构体
type AppUserStruct struct {
	Operation         string `json:"operation"`         //操作
	AppUserID         string `json:"appUserID"`         //应用用户ID
	TxTime            string `json:"txTime"`            //交易时间
	MessageHash       string `json:"messageHash"`       //信息Hash
	Balance           int64  `json:"Balance"`           //积分
	ResearchProjectID string `json:"researchProjectID"` //研究项目ID
	TransactionID     string `json:"transactionID"`     //该交易在链上的TransactionID
}

//BalanceStruct 积分余额结构体
type BalanceStruct struct {
	AppUserID string `json:"AppUserID"` //用户ID
	Balance   int64  `json:"Balance"`   //积分余额
}
