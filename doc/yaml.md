### 格式

+ 大小写敏感
+ 使用缩进表示层级关系
+ 缩进时不允许使用Tab键，只允许使用空格。
+ 缩进的空格数目不重要，只要相同层级的元素左侧对齐即可
+ '#'表示注释，从这个字符一直到行尾，都会被解析器忽略。

### 数据结构

+ 对象：键值对的集合，又称为映射（mapping）/ 哈希（hashes） / 字典（dictionary）
+ 数组：一组按次序排列的值，又称为序列（sequence） / 列表（list）
+ 纯量（scalars）：单个的、不可再分的值

### 对象

animal: pets

hash: { name: Steve, foo: bar } 


### 数组

一组连词线开头的行，构成一个数组。

```
- Cat
- Dog
- Goldfish
```


数据结构的子成员是一个数组，则可以在该项下面缩进一个空格。

```
-
 - Cat
 - Dog
 - Goldfish
```

### 复合结构

```
languages:
 - Ruby
 - Perl
 - Python 
websites:
 YAML: yaml.org 
 Ruby: ruby-lang.org 
 Python: python.org 
 Perl: use.perl.org 
 
 
 { languages: [ 'Ruby', 'Perl', 'Python' ],
  websites: 
   { YAML: 'yaml.org',
     Ruby: 'ruby-lang.org',
     Python: 'python.org',
     Perl: 'use.perl.org' } }
```


### 纯量

```
字符串
布尔值
整数
浮点数
Null
时间
日期
```

null用~表示。


```
number: 12.30
isSet: true
parent: ~ 
```


YAML 允许使用两个感叹号，强制转换数据类型。

```
e: !!str 123
f: !!str true
```

转为 JavaScript 如下

```
{ e: '123', f: 'true' }
```

字符串

```
str: '内容： 字符串'
```