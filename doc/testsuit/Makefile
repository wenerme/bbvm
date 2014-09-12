gbk2utf8=iconv -f gbk -t utf-8

love:
	@echo "Love you too!"

clean:
	find . -name "*.bin" -or -name "*.obj" | xargs rm -rf

list-clean:
	find . -name "*.bin" -or -name "*.obj"

%.bas: .FORCE
	bbasic $@ | $(gbk2utf8)
	

%.obj: %.bas .FORCE
	blink $*.obj $*.bin  | $(gbk2utf8)

%.basm: .FORCE
	blink $*.basm $*.bin  | $(gbk2utf8)

%.obj.r: %.obj
	-cp $*.bin $$BB_HOME/Sim/BBasic/Test.bin;cd $$BB_HOME/Sim/Debug/;./GamDev
	@echo "虚拟结束"

%.run: 
	[ -f $*.bin ] || (echo [ERROR] $*.bin "文件未找到" ; exit 1)
	-cp $*.bin $$BB_HOME/Sim/BBasic/Test.bin;cd $$BB_HOME/Sim/Debug/;./GamDev
	@echo "虚拟结束"

%.exist: .FORCE
	@echo $? - $@ - $^ - $< - $> - $% - $*
	@[ -f $* ] || (echo [ERROR] $* "文件未找到" ; exit 1)

	
abc:
def:
%.tf:
	@echo tf
my:
	echo this is my
%.ft: .FORCE abc def %.tf 
	@echo --------------------
	@echo $? - $@ - $^ - $< - $> - $% - $*

.FORCE:

.PHONY: clean love list-clean

	