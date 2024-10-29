# YU-RIS 引擎 .ybn 文本提取、打包工具

> 基于 [regomne/chinesize/extYbn](https://github.com/regomne/chinesize/tree/master/yuris/extYbn) 修改，主要是增加了对非 ShiftJIS 字符的处理

## Usage

猜测并输出 opcode：

```bash
.\extYbn.exe -e -output-opcode -ybn ".\ysbin\yst00112.ybn"
```

从 ybn 文件提取并输出文本

```bash
.\extYbn.exe -e -v -ybn ".\path\of\yst01234.ybn" ^
  -txt ".\output.txt" -json ".\output.json" ^
  -key 0x96AC6FD3 -op-msg 91 op-call 29
```

写入文本并重新打包 ybn 文件

```bash
.\extYbn.exe -p -v -ybn ".\path\of\yst01234.ybn" ^
  -txt ".\output.txt" -new-ybn "output.ybn" ^
  -key 0x96AC6FD3 -op-msg 91 op-call 29
```
