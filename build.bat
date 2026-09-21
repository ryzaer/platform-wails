@echo off
setlocal EnableExtensions

echo ========================================
echo             APP BUILD
echo ========================================
echo.

set "ROOT=%~dp0"
set "PYTHON=%ROOT%python"
set "INTEGRATOR=%ROOT%integrator"
set "LAUNCHER=%ROOT%launcher"

set "BUILD=%INTEGRATOR%\build\bin"
set "WIN=%BUILD%\windows"
set "LINUX=%BUILD%\linux"

set "PYTHON_DIST=%PYTHON%\renderer.dist"
set "PYTHON_RENDERER=%PYTHON_DIST%\renderer.dll"

set "ISCC=C:\Program Files\Inno Setup 7\ISCC.exe"
set "ISS=%ROOT%launcher\Installer.iss"

echo ROOT       : %ROOT%
echo PYTHON     : %PYTHON%
echo INTEGRATOR : %INTEGRATOR%
echo LAUNCHER   : %LAUNCHER%
echo OUTPUT     : %BUILD%
echo WINDOWS    : %WIN%
echo LINUX      : %LINUX%
echo.

REM ============================================================
REM [1/9] Build Python Renderer
REM ============================================================

echo [1/9] Building Python renderer...
echo.

cd /d "%PYTHON%" || (
    echo ERROR: Cannot enter Python directory.
    exit /b 1
)

if not exist "renderer.py" (
    echo ERROR: renderer.py not found.
    echo Expected: %PYTHON%\renderer.py
    exit /b 1
)

if exist "renderer.build" rmdir /s /q "renderer.build"
if exist "renderer.dist" rmdir /s /q "renderer.dist"

echo Running Nuitka...
python -m nuitka --mode=dll renderer.py

if errorlevel 1 (
    echo.
    echo ERROR: Failed to build Python renderer.
    exit /b 1
)

if not exist "%PYTHON_RENDERER%" (
    echo.
    echo ERROR: renderer.dll was not created.
    echo Expected: %PYTHON_RENDERER%
    exit /b 1
)

echo.
echo Python renderer OK
echo.

REM ============================================================
REM [2/9] Prepare Output Directories
REM ============================================================

echo [2/9] Preparing output directories...
echo.

if exist "%WIN%" rmdir /s /q "%WIN%"
if exist "%LINUX%" rmdir /s /q "%LINUX%"

mkdir "%WIN%"
if errorlevel 1 (
    echo ERROR: Failed to create Windows output directory.
    exit /b 1
)

mkdir "%LINUX%"
if errorlevel 1 (
    echo ERROR: Failed to create Linux output directory.
    exit /b 1
)

echo Output directories OK
echo.

REM ============================================================
REM [4/9] Build Wails integrator.dll
REM ============================================================

echo [4/9] Building Wails integrator.dll...
echo.

cd /d "%INTEGRATOR%" || (
    echo ERROR: Cannot enter integrator directory.
    exit /b 1
)

set "GOOS=windows"
set "GOARCH=amd64"
set "CGO_ENABLED=1"
set "PATH=C:\msys64\ucrt64\bin;%PATH%"

where gcc >nul 2>&1
if errorlevel 1 (
    echo ERROR: GCC not found.
    echo Expected: C:\msys64\ucrt64\bin\gcc.exe
    exit /b 1
)

go build -buildmode=c-shared -tags "desktop,production" -o "integrator.dll" .

if errorlevel 1 (
    echo.
    echo ERROR: Failed to build integrator.dll
    exit /b 1
)

if not exist "%INTEGRATOR%\integrator.dll" (
    echo.
    echo ERROR: integrator.dll was not created.
    exit /b 1
)

echo integrator.dll OK
echo.


REM ============================================================
REM [7/9] Build C++ App Launcher
REM ============================================================

echo [7/9] Building C++ app.exe...
echo.

cd /d "%LAUNCHER%" || (
    echo ERROR: Cannot enter launcher directory.
    exit /b 1
)

set "VCVARS=C:\Program Files\Microsoft Visual Studio\18\Community\VC\Auxiliary\Build\vcvars64.bat"

if not exist "%VCVARS%" (
    echo ERROR: Visual Studio x64 environment not found.
    echo Expected: %VCVARS%
    exit /b 1
)

call "%VCVARS%"
if errorlevel 1 (
    echo ERROR: Failed to initialize Visual Studio x64 environment.
    exit /b 1
)

where cl >nul 2>&1
if errorlevel 1 (
    echo ERROR: cl.exe not found.
    exit /b 1
)

where rc >nul 2>&1
if errorlevel 1 (
    echo ERROR: rc.exe not found.
    exit /b 1
)

if not exist "main.cpp" (
    echo ERROR: main.cpp not found.
    exit /b 1
)

if not exist "launcher.rc" (
    echo ERROR: launcher.rc not found.
    exit /b 1
)

if not exist "launcher.ico" (
    echo ERROR: launcher.ico not found.
    exit /b 1
)

if exist "launcher.res" del /q "launcher.res"
if exist "main.obj" del /q "main.obj"
if exist "app.exe" del /q "app.exe"

rc launcher.rc
if errorlevel 1 (
    echo.
    echo ERROR: Failed to compile launcher resource.
    exit /b 1
)

if not exist "launcher.res" (
    echo ERROR: launcher.res was not created.
    exit /b 1
)
@REM  ini with terminal 
@REM  cl /EHsc /std:c++17 main.cpp launcher.res /Fe:app.exe /link ole32.lib

cl /EHsc /std:c++17 ^
    main.cpp ^
    launcher.res ^
    /Fe:app.exe ^
    /link ^
    /SUBSYSTEM:WINDOWS ^
    ole32.lib ^
    shell32.lib

if errorlevel 1 (
    echo.
    echo ERROR: Failed to build app.exe
    exit /b 1
)

if not exist "%LAUNCHER%\app.exe" (
    echo.
    echo ERROR: app.exe was not created.
    exit /b 1
)

echo app.exe OK
echo.

REM ============================================================
REM [8/9] Assemble Windows Distribution
REM ============================================================

echo [8/9] Assembling Windows distribution...
echo.

move /y "%LAUNCHER%\app.exe" "%WIN%\app.exe" >nul
if errorlevel 1 (
    echo ERROR: Failed to copy app.exe
    exit /b 1
)

move /y "%INTEGRATOR%\integrator.dll" "%WIN%\integrator.dll" >nul
if errorlevel 1 (
    echo ERROR: Failed to copy integrator.dll
    exit /b 1
)

if exist "launcher.res" del /q "launcher.res"
if exist "%INTEGRATOR%\integrator.h" del /q "%INTEGRATOR%\integrator.h"

echo Copying Python renderer distribution...

xcopy "%PYTHON_DIST%\*" "%WIN%\" /E /I /Y /H >nul
if errorlevel 2 (
    echo ERROR: Failed to copy Python renderer distribution.
    exit /b 1
)

echo Windows distribution OK
echo.

REM ============================================================
REM [9/9] Build Installer + Verify
REM ============================================================

echo [9/9] Building Windows installer...
echo.

if not exist "%ISCC%" (
    echo ERROR: Inno Setup compiler not found.
    echo Expected: %ISCC%
    exit /b 1
)

if not exist "%ISS%" (
    echo ERROR: Installer.iss not found.
    echo Expected: %ISS%
    exit /b 1
)

cd /d "%ROOT%" || (
    echo ERROR: Cannot return to project root.
    exit /b 1
)

"%ISCC%" "%ISS%"

if errorlevel 1 (
    echo.
    echo ERROR: Failed to build Windows installer.
    exit /b 1
)

echo App Setup OK
echo.

REM ============================================================
REM VERIFY
REM ============================================================

echo ========================================
echo                VERIFY
echo ========================================
echo.

echo Checking Windows files...

if not exist "%WIN%\app.exe" (
    echo ERROR: app.exe missing.
    exit /b 1
)

if not exist "%WIN%\integrator.dll" (
    echo ERROR: integrator.dll missing.
    exit /b 1
)

if not exist "%WIN%\renderer.dll" (
    echo ERROR: renderer.dll missing.
    exit /b 1
)

if not exist "%WIN%\python311.dll" (
    echo ERROR: python311.dll missing.
    exit /b 1
)

echo Windows files OK
echo.

echo Checking installer...

if not exist "%BUILD%\app-setup.exe" (
    echo ERROR: app-setup.exe missing.
    echo Expected: %BUILD%\app-setup.exe
    exit /b 1
)

echo Installer OK
echo.

echo ========================================
echo             BUILD SUCCESS
echo ========================================
echo.

echo Windows:
echo   %WIN%\app.exe
echo   %WIN%\integrator.dll
echo   %WIN%\runtime.dll
echo   %WIN%\renderer.dll
echo   %WIN%\python311.dll
echo.

echo Linux:
echo   %LINUX%\runtime
echo.

echo Installer:
echo   %BUILD%\app-setup.exe
echo.

echo ========================================

endlocal
exit /b 0
