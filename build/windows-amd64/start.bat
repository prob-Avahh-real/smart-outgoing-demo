@echo off
echo =========================================
echo 启动智慧交通演示系统
echo =========================================
echo.

REM 检查配置文件
if not exist "config\.env" (
    echo 未找到配置文件，使用默认配置...
    set ADMIN_TOKEN=demo_admin_token
    set AMAP_JS_KEY=45109d104b3c8d03a2c84175a7749241
    set AMAP_SECURITY_CODE=c552677838e5f5e71de92ce532c936bc
) else (
    echo 加载配置文件...
    for /f "tokens=*" %%a in ('type config\.env') do set %%a
)

REM 启动服务
echo 启动服务...
smart-traffic.exe
pause
