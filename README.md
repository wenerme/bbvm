BBVM - BeBasic Virtual Machine
===========================
[![Build Status](https://travis-ci.org/wenerme/bbvm.svg)](https://travis-ci.org/wenerme/bbvm)
[![Coverage Status](https://coveralls.io/repos/wenerme/bbvm/badge.svg?branch=master&service=github)](https://coveralls.io/github/wenerme/bbvm?branch=master)
[![Build with love](https://img.shields.io/badge/bbvm-%F0%9F%92%97-orange.svg)](https://github.com/wenerme)
[![GitHub issues](https://img.shields.io/github/issues/wenerme/bbvm.svg)](https://github.com/wenerme/bbvm/issues)
[![GitHub license](https://img.shields.io/badge/license-Apache%202-blue.svg)](https://raw.githubusercontent.com/wenerme/bbvm/master/LICENSE)
[![Twitter](https://img.shields.io/twitter/url/https/github.com/wenerme/bbvm.svg?style=social)](https://twitter.com/intent/tweet?text=Wow:&url=https://github.com/wenerme/bbvm/)

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

第三版实现,第一次实现做到了图像处理,基本实现了功能,但难以扩展.
第二次尝试使用事件处理,使实现过于复杂.
第三次实现尽量保持简洁,按照实际逻辑进行分类,保持 CPU, Memory 等概念,做到可扩展.

目前主要是第四版实现.
完全重写的一个实现,首先基于 JavaCC 实现了解析和编译,使用依赖注入来解耦系统调用和相关的资源管理.
实现简洁,层级简单.测试完全基于自身的解析编译执行,去掉了对原生编译器的依赖.

