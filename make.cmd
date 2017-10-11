@set "args=%*"
@powershell "iex((@('')*3+(cat '%~f0'|select -skip 3))-join[char]10)"
@exit /b %ERRORLEVEL%

Set-PSDebug -strict
$VerbosePreference = "Continue"

function Do-Copy($src,$dst){
    Write-Verbose ('$ copy "{0}" "{1}"' -f $src,$dst)
    Copy-Item $src $dst
}

function Do-Remove($file){
    Write-Verbose ('$ del "{0}"' -f $file)
    Remove-Item $file
}

Add-Type -Assembly System.Windows.Forms
function Ask-Copy($src,$dst){
    $fname = (Join-Path $dst (Split-Path $src -Leaf))
    if( Test-Path $fname ){
        if( "Yes" -ne [System.Windows.Forms.MessageBox]::Show(
            'Override "{0}" by default ?' -f $fname,
            "NYAGOS Install", "YesNo","Question","button2") ){
            return
        }
    }
    Do-Copy $src $dst
}

function Get-GoArch{
    if( Test-Path "goarch.txt" ){
        $arch = (Get-Content "goarch.txt")
    }else{
        $arch = (go version | %{ $_.Split()[-1].Split("/")[-1] } )
    }
    Write-Verbose ("Found GOARCH="+$arch)
    return $arch
}

function ForEach-GoDir{
    Get-ChildItem . -Recurse |
    ?{ $_.Extension -eq '.go' } |
    %{ Split-Path $_.FullName -Parent } |
    Sort-Object |
    Get-Unique
}

function Get-Imports {
    Get-ChildItem . -Recurse |
    ?{ $_.Extension -eq '.go' } |
    %{ Get-Content $_.FullName  } |
    %{ ($_ -replace '\s*//.*$','').Split()[-1] <# remove comment #>} |
    ?{ $_ -match 'github.com/' -and -not ($_ -match '/nyagos/') } |
    %{ $_ -replace '"','' } |
    Sort-Object |
    Get-Unique
}

function Go-Fmt{
    Get-ChildItem . -Recurse |
    ?{ $_.Name -like "*.go" -and $_.Mode -like "?a*" } |
    %{
        $fname = $_.FullName
        Write-Verbose -Message "$ go fmt $fname"
        go fmt $fname
        attrib -a $fname
    }
    Get-ChildItem . -Recurse | ?{ $_.Name -eq "syscall.go" } | %{
        $dir = (Split-Path $_.FullName -Parent)
        $dst = (Join-Path $dir "zsyscall.go")
        if( -not (Test-Path $dst) ){
            Write-Verbose -Message ( `
                "Found {0} but not found $dst" `
                -f $_.FullName,$dst )
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

function Make-SysO($version) {
    Download-Exe "github.com/josephspurrier/goversioninfo/cmd/goversioninfo" "goversioninfo.exe"
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
        ("-ver-patch=" + $v[2]) `
        ("-ver-build=" + $v[3]) `
        ("-product-ver-major=" + $v[0]) `
        ("-product-ver-minor=" + $v[1]) `
        ("-product-ver-patch=" + $v[2]) `
        ("-product-ver-build=" + $v[3]) `
        "-o" nyagos.syso `
        versioninfo.json
}

function Show-Version($fname) {
    if( Test-Path $fname ){
        Write-Output $fname
        $v = [System.Diagnostics.FileVersionInfo]::GetVersionInfo($fname)
        if( $v ){
            Write-Output ("  FileVersion:    `"{0}`" ({1},{2},{3},{4})" -f
                $v.FileVersion,
                $v.FileMajorPart,
                $v.FileMinorPart,
                $v.FileBuildPart,
                $v.FilePrivatePart)
            Write-Output ("  ProductVersion: `"{0}`" ({1},{2},{3},{4})" -f
                $v.ProductVersion,
                $v.ProductMajorPart,
                $v.ProductMinorPart,
                $v.ProductBuildPart,
                $v.ProductPrivatePart)
        }
        $data = [System.IO.File]::ReadAllBytes($fname)
        $md5 = New-Object System.Security.Cryptography.MD5CryptoServiceProvider
        $bs = $md5.ComputeHash($data)
        Write-Output ("  md5sum: {0}" -f [System.BitConverter]::ToString($bs).ToLower().Replace("-",""))
    }else{
        Write-Error ("{0} not found" -f $fname)
    }
}

function Download-Exe($url,$exename){
    if( Test-Path $exename ){
        Write-Verbose -Message ("Found {0}" -f $exename)
        return
    }
    Write-Verbose -Message ("{0} not found." -f $exename)
    Write-Verbose -Message ("$ go get " + $url)
    go get $url
    $workdir = (Join-Path (Join-Path (Get-Go1stPath) "src") $url)
    $cwd = (Get-Location)
    Set-Location $workdir
    Write-Verbose -Message ("$ go build {0} on {1}" -f $exename,$workdir)
    go build
    Do-Copy $exename $cwd
    Set-Location $cwd
}

function Newer-Than($folder,$target){
    if( -not (Test-Path $target) ){
        Write-Verbose ("{0} not found." -f $target)
        return $true
    }
    $stamp = (Get-ItemProperty $target).LastWriteTime

    Get-ChildItem $folder -Recurse | %{
        if( $_.LastWriteTime -gt $stamp ){
            Write-Verbose ("{0} is newer than {1}" -f $_.FullName,$target)
            return $true
        }
    }
    return $false
}


function Build($version,$tags) {
    Write-Verbose -Message ("Build as version='{0}' tags='{1}'" -f $version,$tags)
    Go-Fmt
    $saveGOARCH = $env:GOARCH
    $env:GOARCH = (Get-GoArch)

    Make-SysO $version

    if( Newer-Than "nyagos.d" "mains\bindata.go" ){
        Download-Exe "github.com/jteeuwen/go-bindata/go-bindata" "go-bindata.exe"
        Write-Verbose -Message "$ go-bindata.exe"
        .\go-bindata.exe -pkg "mains" -o "mains\bindata.go" "nyagos.d/..."
    }

    $ldflags = (git log -1 --date=short --pretty=format:"-X main.stamp=%ad -X main.commit=%H")
    Write-Verbose -Message "$ go build"
    go build "-o" nyagos.exe -ldflags "$ldflags -X main.version=$version" $tags
    $env:GOARCH = $saveGOARCH
}

function Make-CSource($package,$names){
    Write-Output '#include <stdio.h>'
    Write-Output '#include <windows.h>'
    Write-Output ''

    $const = @()
    foreach( $p in $names ){
        if( $p -like '"*"' ){
            Write-Output ('#include '+$p)
            continue
        }
        $name1,$type,$fmt = ($p -split ":")
        if( $fmt -ne $null -and $fmt -ne "" ){
            Write-Output `
                ('#define MAKECONST_{0}(n) printf("const " #n "={1}\n",n)' `
                    -f $name1,$fmt)
        }elseif( $type -eq $null -or $type -eq "" ){
            Write-Output `
                ('#define MAKECONST_{0}(n) printf("const " #n "=%d\n",n)' `
                -f $name1)
        }else{
            Write-Output `
                ('#define MAKECONST_{0}(n) printf("const " #n "={1}(%d)\n",n)' `
                -f $name1,$type)
        }
        $const = $const + $name1
    }

    Write-Output 'int main()'
    Write-Output '{'
    Write-Output ('    printf("package {0}\n\n");' -f $package )

    foreach($name1 in $const){
        Write-Output ('    MAKECONST_{0}({0});' -f $name1)
    }
    Write-Output '    return 0;'
    Write-Output '}'
}

function Make-ConstGo($package,$names){
    Make-CSource $package $names |
        Out-File "makeconst.c" -Encoding Default

    Write-Verbose -Message '$ gcc makeconst.c'
    gcc "makeconst.c"

    if( -not (Test-Path "a.exe") ){
        Write-Error -Message "a.exe not found"
        return
    }
    Write-Verbose -Message '$ .\a.exe > const.go'
    & ".\a.exe" | Out-File const.go -Encoding default
    Write-Verbose -Message '$ go fmt const.go'
    go fmt const.go

    if( Test-Path "makeconst.o" ){
        Remove-Item "makeconst.o"
    }
    if( Test-Path "makeconst.c" ){
        Remove-Item "makeconst.c"
    }
    if( Test-Path "a.exe" ){
        Remove-Item "a.exe"
    }
}

$args = $env:args -split " "

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
        ForEach-GoDir | %{
            pushd $_
            Write-Verbose ("$ go vet on " + $_)
            go vet
            popd
        }
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

        ForEach-GoDir | %{
            Write-Verbose -Message ("$ go clean on " + $_)
            pushd $_
            go clean
            popd
        }
    }
    "sweep" {
        Get-ChildItem .  -Recurse |
        ?{ $_.Name -like "*~" -or $_.Name -like "*.bak" } |
        %{ Do-Remove $_.FullName }
    }
    "const" {
        Get-ChildItem . -Recurse |
        ?{ $_.Name -eq "makeconst.txt" } |
        %{
            $private:const = (Get-Content $_.FullName)
            $private:cwd = (Split-Path $_.FullName -Parent)
            pushd $cwd
            Write-Verbose ("$ chdir " + $cwd)
            $private:pkg = (Split-Path $cwd -Leaf)
            Write-Verbose ("for package " + $pkg)
            Make-ConstGo $pkg $const
            popd
        }
    }
    "package" {
        .\nyagos -e "print(nyagos.version or nyagos.stamp)" |
            %{ $version = ($_ -replace "/","") } # get the last line only.
        $private:zipname = ("nyagos-{0}-{1}.zip" -f $version,(Get-GoArch))
        Write-Verbose ("$ zip -9 " + $zipname + " ....")
        zip -9 $zipname `
            nyagos.exe `
            lua53.dll `
            nyagos.lua `
            .nyagos `
            _nyagos `
            makeicon.cmd `
            nyagos.d\*.lua `
            nyagos.d\catalog\*.lua `
            LICENSE `
            readme_ja.md `
            readme.md `
            Doc\*.md
    }
    "install" {
        $installDir = $args[1]
        if( $installDir -eq $null -or $installDir -eq "" ){
            $installDir = (
                Select-String 'INSTALLDIR=([^\)"]+)' Misc\version.cmd |
                %{ $_.Matches.Groups[1].Value }
            )
            if( $installDir -eq $null -or $installDir -eq "" ){
                Write-Warning -Message "Usage: make.ps1 install INSTALLDIR"
                exit
            }
            if( -not (Test-Path $installDir) ){
                Write-Warning -Message ("{0} not found." -f $installDir)
                exit
            }
            Write-Verbose -Message ("installDir="+$installDir)
        }
        Write-Output ('@set "INSTALLDIR={0}"' -f $installDir) |
            Out-File "Misc\version.cmd" -Encoding Default

        robocopy nyagos.d (Join-Path $installDir "nyagos.d") /E
        Write-Verbose ("ERRORLEVEL=" + $LastExitCode)
        if( $LastExitCode -lt 8 ){
            Remove-Item Variable:LastExitCode
        }
        if( -not (Test-Path (Join-Path $installDir "lua53.dll") ) ){
            Do-Copy lua53.dll $installDir
        }
        Do-Copy nyagos.lua $installDir
        Ask-Copy "_nyagos" $installDir
        try{
            Do-Copy nyagos.exe $installDir
        }catch{
            taskkill /F /im nyagos.exe
            Do-Copy nyagos.exe $installDir
            # [void]([System.Windows.Forms.MessageBox]::Show("Done"))
            timeout /T 3
        }
    }
    "get" {
        Get-Imports | ForEach-Object `
            -Process { $_ } `
            -End { "golang.org/x/sys/windows" } |
            %{ Write-Output $_ ; go get -u $_ }
    }
    "fmt" {
        Go-Fmt
    }
    "check-case" {
        $private:dic = @{}
        $private:regex = [regex]"\w+(\-\w+)?"
        $private:done = @{}
        $private:fname = if( $args[1] -ne $null -and $args[1] -ne "" ){
            $args[1]
        }else{
            "make.cmd"
        }
        Get-Content $fname | %{
            $regex.Matches( $_ ) | %{
                $private:one = $_.Value
                $private:key = $one.ToUpper()
                if( $dic.ContainsKey( $key ) ){
                    $private:other = $dic[$key]
                    if( $other -cne $one ){
                        $private:output = ("{0},{1}" -f $one,$other)
                        if( -not $done.ContainsKey($output) ){
                            Write-Output $output
                            $done[ $output ] = $true
                        }
                    }
                }else{
                    $dic.Add($key,$one)
                }
            }
        }
    }
}
if( Test-Path Variable:LastExitCode ){
    exit $LastExitCode
}

# vim:set ft=ps1:
