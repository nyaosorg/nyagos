Set-PSDebug -Strict

$saveEncode = $null
if ([Console]::IsOutputRedirected) {
    $saveEncode = [System.Console]::OutputEncoding
    [System.Console]::OutputEncoding=[System.Text.Encoding]::UTF8
}

Get-ChildItem "release_note*.md" -Recurse | Sort-Object { Format-Hex -InputObject $_.FullName } | ForEach-Object{
    $lang = "(English)"
    if ( $_.FullName -like "*ja*" ) {
        $lang = "(Japanese)"
    }
    $flag = 0
    $section = 0
    Get-Content $_.FullName | ForEach-Object {
        if ( $_ -match "^v?[0-9]+\.[0-9]+\.[0-9]+" ){
            $flag++
            if ( $flag -eq 1 ){
                Write-Host ("`r`n### Changes in {0} in {1}" -f ($_,$lang))
            }
        } elseif ($flag -eq 1 ){
            if ( $_ -eq "" ){
                $section++
            }
            if ( $section -ge 1 ){
                Write-Host $_
            }
        }
    }
}

if ( $saveEncode -ne $null ){
    [System.Console]::OutputEncoding=$saveEncode
}

# gist https://gist.github.com/hymkor/50cd1ed60dc94fe50f12658afcb69cbf
