# Vire

一门极简的编程语言，直接编译为原生程序，支持变量、四则运算和系统调用。

## 快速开始

```bash
git clone https://github.com/vire-lang/vire.git
cd vire
go build -o vire.exe vire.go
```
需要安装 Go 1.21 或更高版本：https://go.dev/dl/

## 示例

创建一个 `test.vire` 文件：

```vire
print(100 + 200)
a = 39
print(a)
syscall(a, 0, 0, 0)
```

运行：

```bash
.\vire.exe test.vire
```

输出：

```
300
39
25048
```

## 语法

- 变量赋值：`a = 123`
- 四则运算：`a + b * c`
- 打印：`print(expr)`
- 系统调用：`syscall(num, arg1, arg2, arg3)`

## 许可证

MIT

## 作者

Vire Lang 团队