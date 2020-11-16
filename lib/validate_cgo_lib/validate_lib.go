package main

//#include <stdlib.h>
import "C"
import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
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
func validate(Cmethod *C.char, Curl *C.char, CbodyData *C.char, CswaggerFilePath *C.char, CerrmsgdefPath *C.char) *C.char {
	//转换数据 CtoGo
	method := C.GoString(Cmethod)
	url := C.GoString(Curl)
	bodyData := C.GoString(CbodyData)
	swaggerFilePath := C.GoString(CswaggerFilePath)
	errmsgdefPath := C.GoString(CerrmsgdefPath)

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
	if err := openapi3filter.ValidateRequest(ctx, requestValidationInput); err != nil {
		msg := err.Error()
		e := &ErrorMessage{
			ErrAPI: route.Path,
		}
		decodeerr := e.Decode(msg, errmsgdefPath)
		return C.CString(decodeerr.Error())
	}
	return C.CString("")
}

//release 用于释放 C.CString() 产生的内存地址，使用时必须在注释中导入 <stdlib.h>
//需要搭配 unsafe 包使用
//export release
func release(p *C.char) {
	C.free(unsafe.Pointer(p))
}

//getErrorFirstLine 获取信息的第一行
func getErrorFirstLine(all string) string {
	firstline := strings.Split(all, "\n")[0]
	return firstline
}

//ErrorMessage 是用于储存从错误信息中解析到的关键内容
type ErrorMessage struct {
	ErrAPI   string      //ErrAPI 出错的接口
	Location string      //Location 出错的位置 <Paremeter><Request body>
	ErrorKey string      //ErrorKey 出错的字段
	Reason   string      //Reason 出错的理由
	Err      interface{} //Err 自定义的错误信息，通过解析获得，主要返回内容
}

//Error 用于ErrorMessage实现error接口
func (e *ErrorMessage) Error() string {
	return fmt.Sprintf("\"error\": \"%v\"", e.Err)
}

//Decode 用于ErrorMessage解析提供的配置文件，获取错误信息
func (e *ErrorMessage) Decode(srcerr string, path string) error {
	//原先的错误信息长度为0则返回nil
	if len(srcerr) == 0 {
		return nil
	}
	//解析错误信息
	decodeerr := e.DecodeSrcError(srcerr)
	if decodeerr != nil {
		return decodeerr
	}
	//fmt.Println(e.ErrAPI, e.Location, e.ErrorKey, e.Reason)

	//解析错误信息配置文件
	errconfig, readErr := ioutil.ReadFile(path)
	if readErr != nil {
		return readErr
	}
	errconfigjson := make(map[string]interface{})
	unmarshalErr := json.Unmarshal(errconfig, &errconfigjson)
	if unmarshalErr != nil {
		return unmarshalErr
	}
	e.Err = errconfigjson[e.Location+e.Reason]
	return e
}

//DecodeSrcError 用于ErrorMessage解析原来的错误信息
func (e *ErrorMessage) DecodeSrcError(srcerr string) error {
	hasAnParamError := regexp.MustCompile(" has an error: ")
	//确定是参数存在错误，进行后续匹配
	if len(hasAnParamError.FindAllString(srcerr, -1)) >= 1 {
		positionMsg := hasAnParamError.Split(srcerr, 2)
		//错误位置匹配
		inParameter := regexp.MustCompile("Parameter '")
		inBody := regexp.MustCompile("Request body")
		inPath := regexp.MustCompile("' in path")
		inQuery := regexp.MustCompile("' in query")
		//Parameter 和 request body公用
		doesntMatchRegExp := regexp.MustCompile("doesn't match the regular expression")
		overLimit := regexp.MustCompile("Number must be")
		//Parameter 错误匹配
		invalidType := make(map[string]*regexp.Regexp)
		keyMissing := make(map[string]*regexp.Regexp)
		invalidType["parameter"] = regexp.MustCompile("invalid syntax")
		keyMissing["parameter"] = regexp.MustCompile("must have a value")
		//Body 错误匹配
		invalidType["body"] = regexp.MustCompile("Field must be set to")
		keyMissing["body"] = regexp.MustCompile("Property '.+' is missing")
		//判断出错的位置path|query
		if len(inParameter.FindAllString(positionMsg[0], -1)) == 1 {
			if inPath.MatchString(positionMsg[0]) {
				e.Location = "Path"
				e.ErrorKey = inPath.ReplaceAllString(inParameter.ReplaceAllString(positionMsg[0], ""), "")
			} else if inQuery.MatchString(positionMsg[0]) {
				e.Location = "Query"
				e.ErrorKey = inQuery.ReplaceAllString(inParameter.ReplaceAllString(positionMsg[0], ""), "")
			} else {
				return errors.New("Error in other location in parameter")
			}
			//判断错误类型
			location := "parameter"
			switch {
			case doesntMatchRegExp.MatchString(positionMsg[1]):
				e.Reason = "MatchRegularExpError"
			case overLimit.MatchString(positionMsg[1]):
				e.Reason = "OverLimit"
			case invalidType[location].MatchString(positionMsg[1]):
				e.Reason = "InvalidType"
			case keyMissing[location].MatchString(positionMsg[1]):
				e.Reason = "KeyMissing"
			default:
				return errors.New("Unconsitered error type")
			}

			//判断出错的位置body
		} else if len(inBody.FindAllString(positionMsg[0], -1)) == 1 {
			e.Location = "Body"
			e.ErrorKey = strings.Split(positionMsg[1], "\"")[1]
			//判断错误类型
			location := "body"
			switch {
			case doesntMatchRegExp.MatchString(positionMsg[1]):
				e.Reason = "MatchRegularExpError"
			case overLimit.MatchString(positionMsg[1]):
				e.Reason = "OverLimit"
			case invalidType[location].MatchString(positionMsg[1]):
				e.Reason = "InvalidType"
			case keyMissing[location].MatchString(positionMsg[1]):
				e.Reason = "KeyMissing"
			default:
				return errors.New("Unconsitered error type")
			}
		} else {
			return errors.New("Source string has more than 1 \"Request body '\" or \"Paramenter '\"")
		}
	} else {
		return errors.New("Source string has more than 1 \" has an error: \"")
	}
	return nil
}
