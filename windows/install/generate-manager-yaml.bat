@echo off
setlocal

set "install_dir=%~1"
set "endpoint=%~2"
set "secret_key=%~3"
set "labels=%~4"

echo %install_dir%
echo %endpoint%
echo %secret_key%
echo %labels%

if "%endpoint%"=="" (
    echo Endpoint not specified; Not writing output yaml
    exit /b 0
)

set "managerfile=%install_dir%\manager.yaml"

echo Writing manager yaml
set "endpointField=endpoint: "%endpoint%""
set "secretField=secret_key: "%secret_key%""
set "labelsField=labels: "%labels%""

echo %endpointField% >"%managerfile%"
if not "%secret_key%"=="" (
    echo %secretField% >>"%managerfile%"
)
if not "%labels%"=="" (
    echo %labelsField% >>"%managerfile%"
)


endlocal
