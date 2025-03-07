@echo off
REM 设置目标系统和架构
set GOOS=linux
set GOARCH=amd64

REM 定义输出文件名称
set OUTPUT_FILE=app

REM 获取当前目录名称作为输出文件名（可选）
for %%F in ("%cd%") do set OUTPUT_FILE=%%~nF

REM 编译当前目录
echo Compiling current directory for %GOOS%/%GOARCH%...
go build -o %OUTPUT_FILE%

REM 检查是否编译成功
if %errorlevel% neq 0 (
    echo Compilation failed!
    exit /b %errorlevel%
)

echo Compilation successful! Output file: %OUTPUT_FILE%