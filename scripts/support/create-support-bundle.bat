@echo off
setlocal enabledelayedexpansion

:: Default directory for logs
set "collector_dir=C:\Program Files\observIQ OpenTelemetry Collector"

:: Check if directory exists, and ask for it if not
if not exist "%collector_dir%" (
    echo Directory %collector_dir% does not exist.
    call :get_dir
)

echo Using directory %collector_dir%
:: Format the date and time for filename
set "datestamp=%DATE:~-4%_%DATE:~4,2%_%DATE:~7,2%"
set "timestamp=%TIME:~0,2%_%TIME:~3,2%_%TIME:~6,2%"
:: replace leading space with zero for hours < 10
set "timestamp=%timestamp: =0%"  

:: Create a new directory to hold the copied files
set "output_dir=Support_Bundle_%datestamp%_%timestamp%"
echo Creating directory %output_dir%

mkdir "%output_dir%"

:: Determine whether to copy only the most recent log
set /p response="Do you want to include only the most recent logs? [Y/n] "
if /I "%response%"=="n" (
    xcopy /Y "%collector_dir%\log\*" "%output_dir%\"
) else (
    xcopy /Y "%collector_dir%\log\collector.log" "%output_dir%\"
)

:: Add the configuration file
xcopy /Y "%collector_dir%\config.yaml" "%output_dir%\"

:: Capture system info
systeminfo > "%output_dir%\systeminfo.txt"

:: Note the location of the copied files
echo Files and system info have been copied to directory %CD%\%output_dir%.

:: End of main script
goto :eof

:: Subroutine to get the directory from the user
:get_dir
set /p collector_dir=Please enter the directory for the collector installation: 
if not exist "%collector_dir%" (
    echo Directory %collector_dir% does not exist.
    exit /b
)
goto :eof
