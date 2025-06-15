@echo off
title 🔁 Air + Go 服务重载运行器
setlocal enabledelayedexpansion

REM 创建 logs 文件夹
if not exist logs mkdir logs

REM 先启动 Air 编译器
echo [Air] 启动热重载监听...
start "Air Watcher" cmd /k "air -c .air.dev.toml"

REM 启动主程序
:loop
cls
echo ----------------------------------------
echo 🚀 启动后端服务: tmp\main.exe
echo ----------------------------------------

REM 如果已运行则结束旧服务
tasklist | findstr /i "main.exe" >nul
if !errorlevel! == 0 (
  echo 🛑 停止旧服务...
  taskkill /f /im main.exe >nul
)

REM 运行新 exe
start "" cmd /c tmp\main.exe

REM 等待用户输入刷新（模拟监听）
echo.
echo 🔄 按任意键重新启动服务，或 Ctrl+C 退出...
pause >nul
goto loop
