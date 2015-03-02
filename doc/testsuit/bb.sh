#! /usr/bin.env bash
BB_HOME=$PWD/bb/
BB_TOOL=$BB_HOME/tool
PATH=$PATH:$BB_HOME/tool

export BB_HOME BB_TOOL

gbk2utf8(){
	iconv -f gbk -t utf-8
}

runtool()
{
	type wine 2>&1 1>/dev/null &&
	{
		if echo $1 | grep $BB_HOME; then
			wine "$1" "$@" | gbk2utf8
		else
			wine "$BB_TOOL/$1" "$@" | gbk2utf8
		fi
		return 0
	}
	"$1" "$@" | gbk2utf8
}

gamdev()
{
	cd $BB_HOME/Sim/Debug/
	runtool $BB_HOME/Sim/Debug/GamDev.exe
	cd -
}

bbtool()
{
	runtool BBTool.exe "$@"
}

blink()
{
	runtool BLink.exe "$@"
}

bbasic()
{
	runtool BBasic.exe "$@"
}

bbr()
{
	cp $1 $BB_HOME/Sim/BBasic/Test.bin
	cd $BB_HOME/Sim/Debug/
	gamdev
	cd -
}

bbc()
{
	bbasic $1 | tee stderr |  grep -q Fail && { cat stderr;echo [bbasic] Compile $1 failed;  return 1;}
	cat stderr
	FN=`basename $1 .bas`
	blink $FN.obj $FN.bin
}
bbcr()
{
	bbasic $1
	FN=`basename $1 .bas`
	blink $FN.obj $FN.bin
	bbr $FN.bin
}

bbar()
{
	FN=`basename $1 .basm`
	blink $1 $FN.bin
	bbr $FN.bin
}
bba()
{
	FN=`basename $1 .basm`
	blink $1 $FN.bin
}


bbhelp()
{
	cat <<HELP
bbr <bin>   run bin file
bbc <bas>   compile bas to bin
bbcr <bas>  compile bas and run
bba <asm>   compile asm to bin
bbar <asm>  compile asm and run
HELP
}

bbcr test.bas