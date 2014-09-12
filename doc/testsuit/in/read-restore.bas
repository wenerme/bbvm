read i
read name$ ' 输出: 1
print i
print name$' 输出: wener

read i
read name$' 输出: 2
print i
print name$' 输出: xiao

restore FDATA
read i
read name$' 输出: 1
print i
print name$' 输出: wener

restore SDATA
read i
read name$' 输出: 3
print i
print name$' 输出: 文儿

read i!
read name$
print i!' 输出: 4.123
print name$' 输出: 笑

FDATA:
data 1, "wener", 2,"xiao"
SDATA:
data 3, "文儿", 4.123,"笑"
