# Define the default directory for logs
$collector_dir = "C:\Program Files\observIQ OpenTelemetry Collector"

# Check if the directory exists
if (!(Test-Path $collector_dir)) {
    Write-Host "Directory $collector_dir does not exist."
    $collector_dir = Read-Host -Prompt "Please enter the directory for logs"
    if (!(Test-Path $collector_dir)) {
        Write-Host "Directory $collector_dir does not exist."
        exit
    }
}

# Create a timestamp for output directory name
$datestamp = Get-Date -Format "yyyy_MM_dd_HH_mm_ss"

# Create a new directory to hold the copied files
$output_dir = "Support_Bundle_$datestamp"
New-Item -ItemType Directory -Force -Path $output_dir

# Determine whether to copy only the most recent log
$response = Read-Host -Prompt "Do you want to include only the most recent logs? [Y/n] "
if ($response -eq "n") {
    Copy-Item "$collector_dir\log\*" -Destination "$output_dir\" -Force
} else {
    Copy-Item "$collector_dir\log\collector.log" -Destination "$output_dir\" -Force
}

# Collector Config
Copy-Item "$collector_dir\config.yaml" -Destination "$output_dir\" -Force

# Capture system info
Get-ComputerInfo | Out-File "$output_dir\systeminfo.txt"

# Compress the files into a zip archive
$zip_filename = "$output_dir.zip"
Compress-Archive -Path "$output_dir\*" -DestinationPath $zip_filename -Force

# Remove the original output directory
Remove-Item -Path $output_dir -Force -Recurse

# Note the location of the copied files
Write-Host "Files and system info have been copied to $(Resolve-Path $zip_filename)."
