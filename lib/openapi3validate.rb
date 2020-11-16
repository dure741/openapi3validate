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
  def self.api_validate(request_method, request_url, request_body_str, swagger_spec_path, errmsgdef_path)
    str, strptr = validate(request_method, request_url, request_body_str, swagger_spec_path, errmsgdef_path)
    release(strptr)
    str
  end
end
#---
# 使用：
# Usage:
# Openapi3Validate::validate(requestMethod,requestURL,requestBodyStr,swaggerSpecPath)
# Openapi3Validate::release(pointer)
# puts Openapi3Validate::apiValidate(requestMethod,requestURL,requestBodyStr,swaggerSpecPath)

#---
# class Test
#     extend Openapi3Validate
#     def initialize
#         @a="a"
#         @b="b"
#         @c="c"
#         @d="d"
#     end
#     def te
#         Openapi3Validate::apiValidate(@a,@b,@c,@d)
#     end
# end
