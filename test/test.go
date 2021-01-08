package main

import "fmt"

//type Student struct {
//	Name string
//}

func main() {
	//interface{} append测试
	//var queryStructs, queryStruct interface{}
	//queryStructs = []Student{}
	//queryStruct = Student{Name: "a"}
	//queryStructs = reflect.Append(reflect.ValueOf(queryStructs), reflect.ValueOf(queryStruct))
	//fmt.Print(queryStructs)    //[{a}]

	//append
	//s:=[]string{"a"}
	//c:=[]string{"b","c","d","e"}
	//s = append(s,c[0],c[1],c[2])
	//fmt.Print(s)

	//切片copy and append  将一个数组所有元素追加到另一个数组后面 == 在一个数组前面追加元素
	//aa := []string{"a", "b"}
	//bb := make([]string, len(aa), (cap(aa))*2)
	//bb=[]string{"cc","dd"}
	////copy(bb, aa)
	//bb= append(bb,aa...)
	//fmt.Println("aa:===", aa)
	//fmt.Println("bb:===", bb)

	//截取数组
	//args := []string{"a", "b", "c", "d", "", "","", "","", "", "sdfsdfsd"}
	//fmt.Println("args:===", args[:len(args)-1])
	//arg := []string{"DataProviderStruct"}
	//perPage := 2
	//fmt.Println("sds====",args[perPage:])
	//for _, value := range args[:len(args)-1] {
	//	if value != "" {
	//		arg = append(arg, value)
	//	} else {
	//		break
	//	}
	//}
	//fmt.Println("arg:=====",arg)

	//test Atoi
	//i := "-11"
	//is, err := strconv.Atoi(i)
	//if err !=nil{
	//	fmt.Println(err)
	//}
	//fmt.Println(is)
	//var his []entity.AppUserStruct
	//hi := entity.AppUserStruct{
	//	TxTime:  "20200521115627",
	//	Balance: 0,
	//}
	//hi2 := entity.AppUserStruct{
	//	TxTime:  "20200521115625",
	//	Balance: 0,
	//}
	//hi3 := entity.AppUserStruct{
	//	TxTime:  "20200521115626",
	//	Balance: 0,
	//}
	//his = append(his, hi)
	//his = append(his, hi2)
	//his = append(his, hi3)
	//
	//history := entity.SortListByTime(his)
	//fmt.Println(history)

	/*
		page , perPage , 预期pages
		1 , 2
	    1,10
	    2,2
	    3,1
	*/
	//分页
	page := 2
	perPage :=3 //每页条数
	var pages int
	aa := []string{"1", "2", "3", "4"}
	total := len(aa)
	startCount := (page - 1) * perPage //分页开始条数
	endCount := startCount + perPage   //分页截止条数
	var empty []string
	if total > endCount {
		aa = aa[startCount:endCount]
	} else if total <= endCount && total > startCount {
		aa = aa[startCount:]
	} else {
		aa = empty
	}
	
	if total%perPage != 0 {
		pages = total/perPage + 1
	} else {
		pages = total / perPage
	}
	fmt.Println("====", aa)

	fmt.Println("total:====", total)
	fmt.Println("pages:====", pages)
	fmt.Println("page:====", page)
	fmt.Println("perPage:====", perPage)

	//
	//total
	//totalPage
	//perPage
	//offsets
	//[]struct{}
}
