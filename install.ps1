Set-PSDebug -Strict

function Install-Nyagos($dir){
    Copy-Item _nyagos  -Destination $dir -PassThru
    Copy-Item nyagos.d -Destination $dir -PassThru -Recurse -ErrorAction SilentlyContinue
    Try {
        Copy-Item nyagos.exe -Destination $dir -PassThru -errorAction stop
    }
    Catch{
        $backup = (Join-Path $dir ("nyagos.exe-" + (Get-Random)))
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
