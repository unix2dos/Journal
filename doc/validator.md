### gopkg.in/go-playground/validator.v8


###tag

tag|desc
---|---
required|不为空, 字段要有, 也不能是默认值
len=10 |长度相等 (int->value string->length)
max=10,min=1|最大值,最小值, (int->value string->length)
eq=10 | 相等 int,string->value
ne=10 | 不等
gt=10,gte=10,lt=10,lte=10 | 大于 int->value, string->length
