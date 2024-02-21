Set-PSDebug -Strict

function Install-Nyagos($dir){
    $random = (Get-Random)
    $new_nyagos_d = (Join-Path $dir "nyagos.d")
    if ( (Test-Path $new_nyagos_d) ){
        Move-Item $new_nyagos_d ($new_nyagos_d + "-" + $random) -PassThru
    }
    Copy-Item nyagos.d -Destination $dir -PassThru -Recurse -ErrorAction SilentlyContinue
    Try {
        Copy-Item nyagos.exe -Destination $dir -PassThru -errorAction stop
    }
    Catch{
        $backup = (Join-Path $dir ("nyagos.exe-" + $random))
        Move-Item (Join-Path $dir "nyagos.exe") $backup -PassThru
        Copy-Item nyagos.exe -Destination $dir -PassThru
    }
}

if ( $args.Length -ge 1 ){
    Install-Nyagos $args[0]
} else {
    Get-Command nyagos.exe | ForEach-Object {
        $target = $_.Source
        $dir = (Split-Path -Parent $target)
        $answer = (Read-Host "Update `"${target}`" ? [y|n]")
        if ( $answer -ieq "y" ){
            Install-Nyagos $dir
        }
    }
}
