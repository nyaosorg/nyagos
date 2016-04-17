# Usage:
#    powershell -ExecutionPolicy RemoteSigned
#       -File makeconst.ps1 PACKAGENAME

$gcc_exe = (Get-Command "gcc").Definition
$bin_dir = [IO.Path]::GetDirectoryName( $gcc_exe )
$gcc_dir = [IO.Path]::GetDirectoryName( $bin_dir )
$inc_dir = [IO.Path]::Combine( $gcc_dir , "include" )

Write-Host("package " + $Args[ 0 ])
Write-Host("")

$keywords = ($Args[1 .. $Args.Length])
Get-ChildItem $inc_dir | %{
    if( $_.Extension -eq ".h" ){
        Get-Content $_.FullName | %{
            foreach( $name in $keywords ){
                if( $_ -match "\b$name\b" ){
                    [Console]::Error.Writeline($_)
                    $value = $_ -replace "(DWORD)","(uint32)"
                    $value = $value.Split()[-1]
                    if( $value -ne "" ){
                        Write-Host("const $name = $value")
                    }
                }
            }
        }
    }
}

# vim:set ft=ps1:
