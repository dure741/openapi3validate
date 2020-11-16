#
# Module Openapi3Validate
# 这个模块封装验证 api 的函数

require 'ffi'
#---
# 这个模块中实现了validate方法
# The function: validate is implemented in this module: Openapi3Validate.
module Test
  extend FFI::Library
  ffi_lib 'lib/validate_cgo_lib/validate_lib.so'
  # validate 不处理指针
  # validate dosn't handle pointer
  attach_function :validate, [:string,:string,:string,:string,:string], :strptr
  # release 处理指针
  # release is used tu handle pointer
  attach_function :release, [:pointer],:void
  # api_validate 中加入了处理指针的功能
  # api_validate can validate and handle pointer
  def self.api_validate(request, swagger_spec_path, errmsgdef_path)
    request_method= request.request_method
    request_url= request.url
    request_body_str= request.body.read
    str, strptr = validate(request_method, request_url, request_body_str, swagger_spec_path, errmsgdef_path)
    release(strptr)
    str
  end
end

(0..1000).each do |i|
    #puts "第#{i}次"
    s = Test::api_validate('POST','http://localhost:4567/user_groups?page_num=x-dd','{
        "usergroup_name": "ddd",
        "user_authorities": "string",
        "authenticate": "string",
        "access_ips": [
            "string"
        ],
        "current_user": "string"
    }',"lib/validate_cgo_lib/swagger.yaml","errmsgdef.json")
    puts s
 end