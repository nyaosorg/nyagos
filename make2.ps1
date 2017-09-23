Set-PSDebug -strict
$VerbosePreference = "Continue"

Function Get-GoArch{
    if( Test-Path "goarch.txt" ){
        Get-Content "goarch.txt"
    }else{
        go version | %{ $_.Split()[-1].Split("/")[-1] }
    }
}

Function EachGoDir{
    Get-ChildItem . -Recurse |
    ?{ $_.Extension -eq '.go' } |
    %{ [System.IO.Path]::GetDirectoryName($_.FullName)} |
    Sort-Object |
    Get-Unique
}

Function Get-Imports {
    Get-ChildItem . -Recurse |
    ?{ $_.Extension -eq '.go' } |
    %{ Get-Content $_.FullName  } |
    %{ ($_ -replace '\s*//.*$','').Split()[-1] <# remove comment #>} |
    ?{ $_ -match 'github.com/' } |
    ?{ -not ($_ -match '/nyagos/') } |
    %{ $_ -replace '"','' } |
    Sort-Object |
    Get-Unique
}

function Is-Modified($path){
    ([System.IO.File]::GetAttributes($path) -band
        [System.IO.FileAttributes]::Archive) -eq
        [System.IO.FileAttributes]::Archive
}

function Get-Modified {
    Get-ChildItem . -Recurse -Name |
    ?{ $_ -like "*.go" } |
    ?{ Is-Modified($_) }
}

function Set-NoModified($path) {
    $attr = [System.IO.File]::GetAttributes($path)
    [System.IO.File]::SetAttributes($path,
        $attr -band -bnot [System.IO.FileAttributes]::Archive)
}

function Go-Fmt{
    Get-Modified | %{
        Write-Verbose -Message "$ go fmt $_"
        go fmt -n $_
        Set-NoModified($_)
    }
    Get-ChildItem . -Recurse | ?{ $_.Name -eq "syscall.go" } | %{
        $dir = (Split-Path $_.FullName -Parent)
        $dst = (Join-Path $dir "zsyscall.go")
        if( -not (Test-Path $dst) ){
            Write-Verbose -Message ( `
                "Found {0} but not found $dst" `
                -f $_.Fullname,$dst )
            pushd $dir
            Write-Verbose -Message ("$ go generate on " + $dir)
            go generate
            popd
        }
    }
}

function Get-Go1stPath {
    if( $env:gopath -ne $null -and $env:gopath -ne "" ){
        $gopath = $env:gopath
    }else{
        $gopath = (Join-Path $env:userprofile "go")
    }
    $gopath.Split(";")[0]
}

function Go-VersionInfo($version) {
    Download-Exe "github.com/josephspurrier/goversioninfo" "goversioninfo.exe" "cmd\goversioninfo"
    if( -not ($version -match "^\d+[\._]\d+[\._]\d+[\._]\d+$") ){
        $version = "0.0.0_0"
    }
    $v = $version.Split("[\._]")

    .\goversioninfo.exe `
        "-file-version=$version" `
        "-product-version=$version" `
        "-icon=mains\nyagos.ico" `
        ("-ver-major=" + $v[0]) `
        ("-ver-minor=" + $v[1]) `
        ("-ver-build=" + $v[2]) `
        ("-ver-patch=" + $v[3]) `
        ("-product-ver-major=" + $v[0]) `
        ("-product-ver-minor=" + $v[1]) `
        ("-product-ver-build=" + $v[2]) `
        ("-product-ver-patch=" + $v[3]) `
        "-o" nyagos.syso `
        versioninfo.json
}

function Show-Version($fname) {
    if( Test-Path $fname ){
        Write-Output $fname
        $v = [System.Diagnostics.FileVersionInfo]::GetVersionInfo($fname)
        if( $v ){
            Write-Output(("  FileVersion:    `"{0}`" ({1},{2},{3},{4})" -f
                $v.FileVersion,
                $v.FileMajorPart,
                $v.FileMinorPart,
                $v.FileBuildPart,
                $v.FilePrivatePart))
            Write-Output(("  ProductVersion: `"{0}`" ({1},{2},{3},{4})" -f
                $v.ProductVersion,
                $v.ProductMajorPart,
                $v.ProductMinorPart,
                $v.ProductBuildPart,
                $v.ProductPrivatePart))
        }
        $data = [System.IO.File]::ReadAllBytes($fname)
        $md5 = New-Object System.Security.Cryptography.MD5CryptoServiceProvider
        $bs = $md5.ComputeHash($data)
        Write-Output ("  md5sum: {0}" -f [System.BitConverter]::ToString($bs).ToLower().Replace("-",""))
    }else{
        Write-Error ("{0} not found" -f $fname)
    }
}

function Download-Exe($url,$exename,$builddir){
    if( Test-Path $exename ){
        Write-Verbose -Message ("Found {0}" -f $exename)
        return
    }
    Write-Verbose -Message ("{0} not found." -f $exename)
    Write-Verbose -Message ("$ go get " + $url)
    go get $url
    $workdir = (Join-Path (Join-Path (Get-Go1stPath) "src") $url)
    if( $builddir -ne $null -and $builddir -ne "" ){
        $workdir = (Join-Path $workdir $builddir)
    }
    $cwd = (Get-Location)
    Set-Location $workdir
    Write-Verbose -Message ("$ go build {0} on {1}" -f $exename,$workdir)
    go build
    Write-Verbose -message ("$ copy {0} {1}" -f $exename,$cwd)
    Copy-Item $exename $cwd
    Set-Location $cwd
}

function Build($version,$tags) {
    Write-Verbose -Message ("Build as version='{0}' tags='{1}'" -f $version,$tags)
    Go-Fmt

    Go-VersionInfo $version

    Get-ChildItem ".\nyagos.d" -Recurse |
    ?{ Is-Modified($_.FullName) } |
    %{
        Write-Verbose -Message ("found {0} modified" -f $_.FullName)
        Set-NoModified($_.FullName)
        if( Test-Path "mains\bindata.go" ){
            Write-Verbose -Message "$ del mains\bindata.go"
            Remove-Item "mains\bindata.go"
        }
    }

    if( -not (Test-Path "mains\bindata.go") ){
        Download-Exe "github.com/jteeuwen/go-bindata" "go-bindata.exe" "go-bindata"
        Write-Verbose -Message "$ go-bindata.exe"
        .\go-bindata.exe -pkg "mains" -o "mains\bindata.go" "nyagos.d/..."
    }

    $ldflags = (git log -1 --date=short --pretty=format:"-X main.stamp=%ad -X main.commit=%H")
    Write-Verbose -Message "$ go build"
    go build "-o" nyagos.exe -ldflags "$ldflags -X main.version=$version" $tags
}

switch( $args[0] ){
    "" {
        Build (git describe --tags) ""
    }
    "debug" {
        Build "" "-tags=debug"
    }
    "release" {
        Build (Get-Content Misc\version.txt) ""
    }
    "status" {
        Show-Version ".\nyagos.exe"
    }
    "vet" {
    }
    "clean" {
        foreach( $p in @(`
            "nyagos.exe",`
            "nyagos.syso",`
            "version.now",`
            "mains\bindata.go",`
            "goversioninfo.exe",`
            "go-bindata.exe" ) )
        {
            if( Test-Path $p ){
                Write-Verbose -Message ("Remove " + $p)
                Remove-Item $p
            }
        }
        Write-Verbose "Search zsyscall.go"
        Get-ChildItem "." -Recurse |
        ?{ $_.Name -eq "zsyscall.go" } |
        %{
            Write-Verbose -Message ("Remove " + $_.FullName)
            Remove-Item $_.FullName
        }

        EachGoDir | %{
            Write-Verbose -Message ("$ go clean on " + $_)
            pushd $_
            go clean
            popd
        }
    }
    "sweep" {
    }
    "bindata" {
    }
    "goversioninfo" {
    }
    "const" {
    }
    "package" {
    }
    "install" {
    }
    "get" {
        Get-Imports | %{ Write-Output $_ ; go get -u $_ }
    }
    "fmt" {
        Go-Fmt
    }
}
