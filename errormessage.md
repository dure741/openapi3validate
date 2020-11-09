# 错误信息：
## 参数
### query
---
- 不匹配正则表达式
` Parameter` 'page_num' in ` query` has an error: JSON string doesn't match the ` regular expression` '^x-'
- 不合规值
` Parameter` 'page_num' in ` query` has an error: value ssd: an invalid integer: strconv.ParseFloat: parsing "ssd": invalid syntax
- 数值超限
` Parameter` 'page_size' in ` query` has an error: ` Number must be most` 23
` Parameter` 'page_size' in ` query` has an error: ` Number must be at least` 10
- 缺失查询参数
` Parameter` 'page_num' in ` query` has an error: ` must have a value`: must have a value
### path
- 不合规值
` Parameter` 'id' in ` path` has an error: value jfije: an invalid integer: strconv.ParseFloat: parsing "jfije": invalid syntax
- 数值超限
` Parameter` 'id' in ` path` has an error: ` Number must be at least` 10
### header
### cookie
## 请求体
### body
---
- 键值缺失
` Request body` has an error: doesn't match the schema: Error at "/usergroup_name":` Property` 'usergroup_name' ` is missing`
- 不匹配正则表达式
` Request body` has an error: doesn't match the schema: Error at "/usergroup_name":JSON string doesn't match the ` regular expression` '^x-'
- 值不为字符串
` Request body` has an error: doesn't match the schema: Error at "/current_user":` Field must be set to` string or not be present
- 值不为integer
` Request body` has an error: doesn't match the schema: Error at "/current_user":` Field must be set to` integer or not be present
- 数值超限
` Request body` has an error: doesn't match the schema: Error at "/current_user":` Number must be at least` 10
` Request body` has an error: doesn't match the schema: Error at "/current_user":` Number must be most` 23
## 认证策略
暂调研深入