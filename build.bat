@echo off
echo 正在编译 Vire 引擎...
gcc -o vire.exe src\vire.c -lm
if %errorlevel% == 0 (
    echo ? 编译成功！试试运行: vire.exe examples\hello.vire
) else (
    echo ? 编译失败，请确保已安装 MinGW (gcc)
)
pause
