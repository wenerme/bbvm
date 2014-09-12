s$ = left$("abcdefg", 3)
print "should-be:abc"
print s$

s$ = right$("abcdefg", 3)
print "should-be:efg"
print s$

print "should-be:0"
print INSTR(0,"=x","=x2007 = 2008 - 1")
print "should-be:-1"
print INSTR(0,"=x","=2007 = 2008 - 1")
