@set args=%*
@powershell "iex((@('')*3+(cat '%~f0'|select -skip 3))-join[char]10)"
@exit /b %ERRORLEVEL%

$args = @( ([regex]'"([^"]*)"').Replace($env:args,{
        $args[0].Groups[1] -replace " ",[char]1
    }) -split " " | ForEach-Object{ $_ -replace [char]1," " })

Set-PSDebug -strict
$VerbosePreference = "Continue"
$env:GO111MODULE="on"
Write-Verbose "$ set GO111MODULE=$env:GO111MODULE"

function Do-Copy($src,$dst){
    Write-Verbose "$ copy '$src' '$dst'"
    Copy-Item $src $dst
}

function Do-Rename($old,$new){
    Write-Verbose "$ ren $old $new"
    Rename-Item -path $old -newname $new
}

function Do-Remove($file){
    if( Test-Path $file ){
        Write-Verbose "$ del '$file'"
        Remove-Item $file
    }
}

function Make-Dir($folder){
    if( -not (Test-Path $folder) ){
        Write-Verbose "$ mkdir '$folder'"
        New-Item $folder -type directory | Out-Null
    }
}

function Ask-Copy($src,$dst){
    $fname = (Join-Path $dst (Split-Path $src -Leaf))
    if( Test-Path $fname ){
        Write-Verbose "$fname already exists. Cancel to copy"
    }else{
        Do-Copy $src $dst
    }
}

function ForEach-GoDir{
    Get-ChildItem . -Recurse |
    Where-Object{ $_.Extension -eq '.go' } |
    ForEach-Object{ Split-Path $_.FullName -Parent } |
    Sort-Object |
    Get-Unique
}

function Go-Fmt{
    $status = $true
    git status -s | %{
        $fname = $_.Substring(3)
        $arrow = $fname.IndexOf(" -> ")
        if( $arrow -ge 0 ){
            $fname = $fname.Substring($arrow+4)
        }
        if( $fname -like "*.go" -and (Test-Path($fname)) ){
            $prop = Get-ItemProperty($fname)
            if( $prop.Mode -like "?a*" ){
                Write-Verbose "$ go fmt $fname"
                go fmt $fname
                if( $LastExitCode -ne 0 ){
                    $status = $false
                }else{
                    attrib -a $fname
                }
            }
        }
    }
    if( -not $status ){
        Write-Warning "Some of 'go fmt' failed."
    }
    return $status
}

function Make-SysO($version) {
    if( $env:GOOS -eq $null -or $env:GOOS -eq "" -or $env:GOOS -eq "windows" ){
        pushd Etc
        go generate
        popd
    }
}

function Build([string]$version="",[string]$tags="",[string]$target="") {
    if( $version -eq "" ){
        $version = (git describe --tags)
    }

    Write-Verbose "Build as version='$version' tags='$tags'"

    if( $tags -ne "" ){
        $tags = "-tags=$tags"
    }

    if( -not (Go-Fmt) ){
        return
    }
    $saveGOARCH = $env:GOARCH
    $env:GOARCH = (go env GOARCH)

    Make-Dir "bin"
    $binDir = (Join-Path "bin" $env:GOARCH)
    Make-Dir $binDir
    if ($target -eq "") {
        $target = (Join-Path $binDir "nyagos.exe")
    }

    Make-SysO $version

    Write-Verbose "$ go build -o '$target'"
    go build "-o" $target -ldflags "-s -w -X main.version=$version" $tags
    if( $LastExitCode -eq 0 ){
        Do-Copy $target (Join-Path "." ([System.IO.Path]::GetFileName($target)))
    }
    $env:GOARCH = $saveGOARCH
}

function Make-Package($arch){
    $zipname = ("nyagos-{0}.zip" -f (& "bin\$arch\nyagos.exe" --show-version-only))

    where.exe upx 2>&1 | Out-Null
    if ( $LastExitCode -eq 0 ){
        upx.exe -9 "bin\$arch\nyagos.exe"
    }else{
        $global:LastExitCode = 0
    }

    Write-Verbose "$ zip -9 $zipname ...."
    if( Test-Path $zipname ){
        Do-Remove $zipname
    }

    zip -9j $zipname `
        "bin\$arch\nyagos.exe" `
        .nyagos `
        _nyagos `
        makeicon.cmd `
        LICENSE `
        readme_ja.md `
        readme.md

    zip -9 $zipname `
        nyagos.d\*.lua `
        nyagos.d\catalog\*.lua `
        Doc\*.md
}

switch( $args[0] ){
    "" {
        Build
    }
    "386"{
        $private:save = $env:GOARCH
        $env:GOARCH = "386"
        Build
        $env:GOARCH = $save
    }
    "debug" {
        $private:save = $env:GOARCH
        if( $args[1] ){
            $env:GOARCH = $args[1]
        }
        Build -tags "debug"
        $env:GOARCH = $save
    }
    "vanilla" {
        Build -tags "vanilla"
    }
    "release" {
        $private:save = $env:GOARCH
        if( $args[1] ){
            $env:GOARCH = $args[1]
        }
        Build -version (Get-Content Etc\version.txt)
        $env:GOARCH = $save
    }
    "linux" {
        $private:os = $env:GOOS
        $private:arch = $env:GOARCH
        $env:GOOS="linux"
        $env:GOARCH="amd64"
        Build -target "bin\linux\nyagos" -version (Get-Content Etc\version.txt)
        $env:GOOS = $os
        $env:GOARCH=$arch
    }
    "clean" {
        foreach( $p in @(`
            "bin\amd64\nyagos.exe",`
            "bin\386\nyagos.exe",`
            "nyagos.exe",`
            "nyagos.syso",`
            "version.now",`
            "goversioninfo.exe") )
        {
            Do-Remove $p
        }
        ForEach-GoDir | %{
            Write-Verbose "$ go clean on $_"
            pushd $_
            go clean
            popd
        }
    }
    "package" {
        $ARCH = (go env GOARCH)
        $private:VER = (Get-Content Etc\version.txt)
        if( $args[1] -eq "linux" ){
            pushd ..
            tar -zcvf "nyagos/nyagos-$VER-linux-$ARCH.tar.gz" `
                nyagos/nyagos `
                nyagos/.nyagos `
                nyagos/_nyagos `
                nyagos/readme.md `
                nyagos/readme_ja.md `
                nyagos/nyagos.d `
                nyagos/Doc/*.md
            popd
        }elseif( $args[1] -ne $null -and $args[1] -ne "" ){
            Make-Package $args[1]
        }else{
            Make-Package $ARCH
        }
    }
    "install" {
        $installDir = $args[1]
        if( $installDir -eq $null -or $installDir -eq "" ){
            $installDir = (
                Select-String 'INSTALLDIR=([^\)"]+)' Etc\version.cmd |
                ForEach-Object{ $_.Matches[0].Groups[1].Value }
            )
            if( -not $installDir ){
                Write-Warning "Usage: make.ps1 install INSTALLDIR"
                exit
            }
            if( -not (Test-Path $installDir) ){
                Write-Warning "$installDir not found."
                exit
            }
            Write-Verbose "installDir=$installDir"
        }
        Write-Output "@set `"INSTALLDIR=$installDir`"" |
            Out-File "Etc\version.cmd" -Encoding Default

        robocopy nyagos.d (Join-Path $installDir "nyagos.d") /E
        Write-Verbose "ERRORLEVEL=$LastExitCode"
        if( $LastExitCode -lt 8 ){
            Remove-Item Variable:LastExitCode
        }
        Ask-Copy "_nyagos" $installDir
        try{
            Do-Copy nyagos.exe $installDir
        }catch{
            $old = (Join-Path $installDir "nyagos.exe")
            Write-Warning "Failed to update $old"
            $now = (Get-Date -Format "yyyyMMddHHmmss")
            $backup = ($old + "-" + $now)
            try{
                try{
                    Do-Rename $old $backup
                }catch{
                    Write-Warning "Failed to rename $old to $backup"
                    Write-Warning "Try to kill nyagos.exe process"
                    taskkill /F /IM nyagos.exe
                    Do-Rename $old ($old + "-" + $now)
                }
                Do-Copy nyagos.exe $installDir
            }catch{
                Write-Error "Could not update installed nyagos.exe"
                Write-Error "Some processes holds nyagos.exe now"
            }
        }
    }
    "get" {
        go get -u
    }
    "fmt" {
        Go-Fmt | Out-Null
    }
    "help" {
        Write-Output @'
make                     build as snapshot
make debug   [386|amd64] build as debug version     (tagged as `debug`)
make release [386|amd64] build as release version
make clean               remove all work files
make package [386|amd64] make `nyagos-(VERSION)-(ARCH).zip`
make install [FOLDER]    copy executables to FOLDER or last folder
make fmt                 `go fmt`
make help                show this
'@
    }
    default {
        Write-Warning ("{0} not supported." -f $args[0])
    }
}
if( Test-Path Variable:LastExitCode ){
    exit $LastExitCode
}

# vim:set ft=ps1:
