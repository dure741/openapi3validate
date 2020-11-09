package main

//#include <stdlib.h>
import "C"
import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"unsafe"

	"github.com/getkin/kin-openapi/openapi3filter"
)

func main() {
	router := openapi3filter.NewRouter().WithSwaggerFromFile("swagger.yaml")
	ctx := context.Background()
	httpReq, _ := http.NewRequest(http.MethodDelete, "http://localhost:4567/user_groups/jfjfj", nil)

	// Find route
	route, pathParams, _ := router.FindRoute(httpReq.Method, httpReq.URL)

	// Validate request
	requestValidationInput := &openapi3filter.RequestValidationInput{
		Request:    httpReq,
		PathParams: pathParams,
		Route:      route,
	}
	if err := openapi3filter.ValidateRequest(ctx, requestValidationInput); err != nil {
		fmt.Println(err.Error())
	}

	// var (
	// 	respStatus      = 200
	// 	respContentType = "application/json"
	// 	respBody        = bytes.NewBufferString(`{}`)
	// )

	// log.Println("Response:", respStatus)
	// responseValidationInput := &openapi3filter.ResponseValidationInput{
	// 	RequestValidationInput: requestValidationInput,
	// 	Status:                 respStatus,
	// 	Header: http.Header{
	// 		"Content-Type": []string{respContentType},
	// 	},
	// }
	// if respBody != nil {
	// 	data, _ := json.Marshal(respBody)
	// 	responseValidationInput.SetBodyBytes(data)
	// }

	// Validate response.
	// if err := openapi3filter.ValidateResponse(ctx, responseValidationInput); err != nil {
	// 	panic(err)
	// }
}

//validate 用来验证请求是否符合接口规范的定义
//因为使用Cgo，暂不知如何传入http请求，所以使用字符串来处理，必要的几个参数为
//请求方法 method，请求路径 path，路径查询参数 params，请求体内容 bodyData(暂时只考虑json)，openapi规范文件路径 swaggerFilePath
//因为Cgo，所以返回无法方便使用error，所以使用字符串，
//返回 "" 表示没有验证到错误，返回长度非0的字符串则表示未通过验证出现了错误
//export validate
func validate(Cmethod *C.char, Curl *C.char, CbodyData *C.char, CswaggerFilePath *C.char) *C.char {
	//转换数据 CtoGo
	method := C.GoString(Cmethod)
	url := C.GoString(Curl)
	bodyData := C.GoString(CbodyData)
	swaggerFilePath := C.GoString(CswaggerFilePath)

	router := openapi3filter.NewRouter().WithSwaggerFromFile(swaggerFilePath)
	ctx := context.Background()
	body := strings.NewReader(bodyData)
	httpReq, _ := http.NewRequest(method, url, body)
	//以下修改请求头的语句是写死的，因现在只有application/json
	httpReq.Header.Set("Content-Type", "application/json")
	// Find route
	route, pathParams, _ := router.FindRoute(httpReq.Method, httpReq.URL)
	// Validate request
	requestValidationInput := &openapi3filter.RequestValidationInput{
		Request:    httpReq,
		PathParams: pathParams,
		Route:      route,
	}
	err := openapi3filter.ValidateRequest(ctx, requestValidationInput)
	if err != nil {
		return C.CString(err.Error())
	}
	return C.CString("")
}

//release 用于释放 C.CString() 产生的内存地址，使用时必须在注释中导入 <stdlib.h>
//需要搭配 unsafe 包使用
//export release
func release(p *C.char) {
	C.free(unsafe.Pointer(p))
}
