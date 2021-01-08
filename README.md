//ssh下载项目
git clone ssh://git@nvxg.nvxclouds.com:11622/blockchain/fabric1.4/nvxclouds_chaincode.git

//数据节点注册
peer chaincode invoke -n tobcc -c '{"Args":["dataProviderInvoke","{\"Operation\":\"dataNodeRegistration\",\"DataNodeID\":\"dataNode1\",\"NodeName\":\"111\",\"TxTime\":\"20200521115623\",\"MessageHash\":\"node111\"}"]}' -C nvxchannel --tls true --cafile $ORDERER_CA 
返回结果：
2020-05-21 07:44:38.627 UTC [chaincodeCmd] chaincodeInvokeOrQuery -> INFO 004 Chaincode invoke successful. result: status:200 payload:"{\"TransactionID\":\"58c14b33db511a22f71afaa7abeffb2e51724f8d5163ab29c4282a6e5bf691d0\"}"

peer chaincode invoke -n tobcc -c '{"Args":["dataProviderInvoke","{\"Operation\":\"dataNodeRegistration\",\"DataNodeID\":\"dataNode2\",\"NodeName\":\"222\",\"TxTime\":\"20200521115625\",\"MessageHash\":\"node222\"}"]}' -C nvxchannel --tls true --cafile $ORDERER_CA

peer chaincode invoke -n tobcc -c '{"Args":["dataProviderInvoke","{\"Operation\":\"dataNodeRegistration\",\"DataNodeID\":\"dataNode3\",\"NodeName\":\"333\",\"TxTime\":\"20200521115626\",\"MessageHash\":\"node333\"}"]}' -C nvxchannel --tls true --cafile $ORDERER_CA

//精准查询数据节点注册信息
peer chaincode query -n tobcc -c '{"Args":["dataProviderQuery","{\"Operation\":\"dataNodeRegistration\",\"DataNodeID\":\"dataNode1\",\"NodeName\":\"111\",\"TxTime\":\"20200521115623\",\"MessageHash\":\"node111\"}"]}' -C nvxchannel
返回结果：
[{"operation":"dataNodeRegistration","dataNodeID":"dataNode1","txTime":"20200521115623","nodeName":"111","messageHash":"node111","datasetID":"","status":"","transactionID":"58c14b33db511a22f71afaa7abeffb2e51724f8d5163ab29c4282a6e5bf691d0"}]

//范围查询(前缀查询)
//根据操作名称，数据节点ID查询数据节点注册信息
peer chaincode query -n tobcc -c '{"Args":["dataProviderQuery","{\"Operation\":\"dataNodeRegistration\",\"DataNodeID\":\"dataNode1\"}"]}' -C nvxchannel

//根据操作名称，数据节点ID，节点名称查询数据节点注册信息
peer chaincode query -n tobcc -c '{"Args":["dataProviderQuery","{\"Operation\":\"dataNodeRegistration\",\"DataNodeID\":\"dataNode1\",\"NodeName\":\"111\"}"]}' -C nvxchannel

//根据操作名称，数据节点ID，节点名称，注册时间查询数据节点注册信息
peer chaincode query -n tobcc -c '{"Args":["dataProviderQuery","{\"Operation\":\"dataNodeRegistration\",\"DataNodeID\":\"dataNode1\",\"NodeName\":\"111\",\"TxTime\":\"20200521115623\"}"]}' -C nvxchannel

//根据操作名称，数据节点ID，起止时间范围查询
peer chaincode query -n tobcc -c '{"Args":["dataProviderQuery","{\"Operation\":\"dataNodeRegistration\",\"DataNodeID\":\"dataNode1\",\"StartTime\":\"20200521115622\",\"EndTime\":\"20200521115624\"}"]}' -C nvxchannel
返回结果：
[{"operation":"dataNodeRegistration","dataNodeID":"dataNode1","txTime":"20200521115623","nodeName":"111","messageHash":"node111","datasetID":"","status":"","transactionID":"58c14b33db511a22f71afaa7abeffb2e51724f8d5163ab29c4282a6e5bf691d0"}]

//根据操作名称，起止时间范围查询
peer chaincode query -n tobcc -c '{"Args":["dataProviderQuery","{\"Operation\":\"dataNodeRegistration\",\"StartTime\":\"20200521115622\",\"EndTime\":\"20200521115627\"}"]}' -C nvxchannel
返回结果：
[{"operation":"dataNodeRegistration","dataNodeID":"dataNode1","txTime":"20200521115623","nodeName":"111","messageHash":"node111","datasetID":"","status":"","transactionID":"58c14b33db511a22f71afaa7abeffb2e51724f8d5163ab29c4282a6e5bf691d0"},{"operation":"dataNodeRegistration","dataNodeID":"dataNode2","txTime":"20200521115625","nodeName":"222","messageHash":"node222","datasetID":"","status":"","transactionID":"cfe65f06da1dad5764839dc019854d225fc9a220d7fdfe3b221b34e7badebbc8"},{"operation":"dataNodeRegistration","dataNodeID":"dataNode3","txTime":"20200521115626","nodeName":"333","messageHash":"node333","datasetID":"","status":"","transactionID":"86991bf0003db382c86e4502f13d719db5ef63bab0c8c57fb9a1ac1d39403ca5"}]

//根据操作名称，起止时间,跳过页数范围查询
peer chaincode query -n tobcc -c '{"Args":["dataProviderQuery","{\"Operation\":\"dataNodeRegistration\",\"StartTime\":\"20200521115622\",\"EndTime\":\"20200521115627\",\"PageSize\":\"1\"}"]}' -C nvxchannel
返回结果：
[{"operation":"dataNodeRegistration","dataNodeID":"dataNode2","txTime":"20200521115625","nodeName":"222","messageHash":"node222","datasetID":"","status":"","transactionID":"cfe65f06da1dad5764839dc019854d225fc9a220d7fdfe3b221b34e7badebbc8"},{"operation":"dataNodeRegistration","dataNodeID":"dataNode3","txTime":"20200521115626","nodeName":"333","messageHash":"node333","datasetID":"","status":"","transactionID":"86991bf0003db382c86e4502f13d719db5ef63bab0c8c57fb9a1ac1d39403ca5"}]