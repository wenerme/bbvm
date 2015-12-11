BBvm - BeBasic Virtual Machine
===========================

该项目为原步步高BBasic的一个仿照实现.并在原来的基础上进行了扩展.

主要目标

* 做到和 BBasic 的汇编码兼容
* 做到和 BBasic 的二进制兼容
* 实现编译 BBAsm 的编译器
* 实现 BB 的虚拟机,包括图形界面等所有功能
* 对 BBAsm 进行扩展

参考
====

* BBAsm 语法参考[这里][bbasm-g4]
* BB 虚拟机规范参考[这里][bbvm-spec]

 [bbasm-g4]:https://github.com/wenerme/bbvm/blob/master/doc/grammar/BBAsm.g4
 [bbvm-spec]:https://github.com/wenerme/bbvm/wiki/vm-spec



NOTE
====

目前是第三版实现,第一次实现做到了图像处理,基本实现了功能,但难以扩展.
第二次尝试使用事件处理,使实现过于复杂.
第三次实现尽量保持简洁,按照实际逻辑进行分类,保持 CPU, Memory 等概念,做到可扩展.

目前主要是第三版实现.
