BBvm - BeBasic Virtual Machine
===========================

该项目为原步步高BBasic的一个仿照实现.并在原来的基础上进行了扩展.

主要目标
------
* 做到和 BBasic 的汇编码兼容
* 做到和 BBasic 的二进制兼容
* 实现编译 BBAsm 的编译器
* 实现 BB 的虚拟机,包括图形界面等所有功能
* 对 BBAsm 进行扩展
* 图片格式支持

TODO
----
* 图片格式支持
	* [x] RLB 解码
	* [ ] RLB 编码
	* [X] LIB BE/LE 解码
	* [ ] LIB BE/LE 编码
	* [ ] DLX 解码
	* [ ] DLX 编码
* 支持图像端口操作
	* [x] 页面操作 CREATEPAGE,DELETEPAGE,FLIPPAGE
	* [X] 页面绘图 CIRCLE,RECTANGE,PIXEL,READPIXEL
	* [ ] 页面文字输出
	* [ ] 字体选择端口
	* [ ] 资源操作 LOADRES,DELETERES
	* [ ] 显示图片资源
	* [ ] 页面互操作
* 使用 SDL 输出图像
* 编译器
	
动机
===
* 祭奠曾经
* 学习 Go
* 完成自己的梦想

参考
====

* BBAsm 语法参考[这里][bbasm-g4]
* BB 虚拟机规范参考[这里][bbvm-spec]
* [Java 版 BBVM](https://github.com/wenerme/bbvm/tree/java)

 [bbasm-g4]:https://github.com/wenerme/bbvm/blob/master/doc/grammar/BBAsm.g4
 [bbvm-spec]:https://github.com/wenerme/bbvm/blob/master/doc/bbvm-spec.md


