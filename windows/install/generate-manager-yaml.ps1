param (
    [string] $install_dir = $null,
    [string] $endpoint = $null,
    [string] $secret_key = $null,
    [string] $labels = $null
)

Write-Output $install_dir
Write-Output $endpoint
Write-Output $secret_key
Write-Output $labels

if([string]::IsNullOrEmpty($endpoint)) {
    # Don't generate yaml if endpoint isn't specified
    Write-Output "Endpoint not specified; Not writing output yaml"
    return 0
}

Write-Output "Writing manager yaml"
$yaml="endpoint: `"$endpoint`"{0}" -f [environment]::NewLine
if(![string]::IsNullOrEmpty($secret_key)) {
    $yaml+="secret_key: `"$secret_key`"{0}" -f [environment]::NewLine
}
if(![string]::IsNullOrEmpty($labels)) {
    $yaml+="labels: `"$labels`"{0}" -f [environment]::NewLine
}
$yaml+="agent_id: `"{0}`"" -f [guid]::NewGuid()

Set-Content -Path "${install_dir}manager.yaml" -Value $yaml
