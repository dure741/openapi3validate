#
# Module Openapi3Validate
# 这个模块封装验证 api 的函数

require 'ffi'
#---
# 这个模块中实现了validate方法
# The function: validate is implemented in this module: Openapi3Validate.
module Openapi3Validate
  extend FFI::Library
  ffi_lib '/home/durui/.gem/ruby/gems/openapi3validate-1.0.0/lib/validate_cgo_lib/validate_lib.so'
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

