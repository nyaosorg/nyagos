@set args=%*
@powershell "iex((@('')*3+(cat '%~f0'|select -skip 3))-join[char]10)"
@exit /b %ERRORLEVEL%

$args = @( ([regex]'"([^"]*)"').Replace($env:args,{
        $args[0].Groups[1] -replace " ",[char]1
    }) -split " " | ForEach-Object{ $_ -replace [char]1," " })

# set GO "go" -option constant
set GO "go.exe" -option constant

set CMD "Cmd" -option constant

$LUAURL = @{
    "amd64"="https://sourceforge.net/projects/luabinaries/files/5.3.4/Windows%20Libraries/Dynamic/lua-5.3.4_Win64_dllw4_lib.zip/download";
    "386"="https://sourceforge.net/projects/luabinaries/files/5.3.4/Windows%20Libraries/Dynamic/lua-5.3.4_Win32_dllw4_lib.zip/download";
}

Set-PSDebug -strict
$VerbosePreference = "Continue"

function Do-Copy($src,$dst){
    Write-Verbose "$ copy '$src' '$dst'"
    Copy-Item $src $dst
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
    }elseif( (Test-Path env:GOARCH) -and $env:GOARCH ){
        $arch = $env:GOARCH
    }else{
        $arch = (& $GO version | %{ $_.Split()[-1].Split("/")[-1] } )
    }
    Write-Verbose ("Found GOARCH="+$arch)
    return $arch
}

function ForEach-GoDir{
    Get-ChildItem . -Recurse |
    Where-Object{ $_.Extension -eq '.go' } |
    ForEach-Object{ Split-Path $_.FullName -Parent } |
    Sort-Object |
    Get-Unique
}

function Get-Imports {
    Get-ChildItem . -Recurse |
    Where-Object{ $_.Extension -eq '.go' } |
    ForEach-Object{ Get-Content $_.FullName  } |
    ForEach-Object{ ($_ -replace '\s*//.*$','').Split()[-1] <# remove comment #>} |
    ?{ ($_ -match 'github.com/' -and $_ -notmatch '/nyagos/') `
        -or $_ -match 'golang.org/' } |
    %{ $_ -replace '"','' } |
    Sort-Object |
    Get-Unique
}

function Go-Generate{
    Get-ChildItem "." -Recurse |
    Where-Object{ $_.Name -eq "make.xml" } |
    ForEach-Object{
        $dir = (Split-Path $_.FullName -Parent)
        pushd $dir
        $xml = [xml](Get-Content $_.FullName)
        :allloop foreach( $li in $xml.make.generate.li ){
            foreach( $target in $li.target ){
                if( -not $target ){ continue }
                foreach( $source in $li.source ){
                    if( -not $source ){ continue }
                    if( (Newer-Than $source $target) ){
                        Write-Verbose ("$ $GO generate for {0}" -f
                            (Join-Path $dir $target) )
                        & $GO generate
                        break allloop
                    }
                }
            }
        }
        popd
    }
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
                Write-Verbose "$ $GO fmt $fname"
                & $GO fmt $fname
                if( $LastExitCode -ne 0 ){
                    $status = $false
                }else{
                    attrib -a $fname
                }
            }
        }
    }
    if( -not $status ){
        Write-Warning "Some of '$GO fmt' failed."
    }
    return $status
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
    if( $version -match "^\d+[\._]\d+[\._]\d+[\._]\d+$" ){
        $v = $version.Split("[\._]")
    }else{
        $v = @(0,0,0,0)
        if( $version -eq $null -or $version -eq "" ){
            $version = "0.0.0_0"
        }
    }
    Write-Verbose "version=$version"

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
    if( -not (Test-Path $fname) ){
        Write-Error ("{0} not found" -f $fname)
        return
    }
    Write-Output $fname
    $data = [System.IO.File]::ReadAllBytes($fname)
    $bits = switch( (Get-Architecture $data) ){
        32 { "32bit or AnyCPU" }
        64 { "64bit" }
        $null { "unknown"}
    }
    Write-Output ("  Architecture:   {0}" -f $bits)
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
    $md5 = New-Object System.Security.Cryptography.MD5CryptoServiceProvider
    $bs = $md5.ComputeHash($data)
    Write-Output ("  md5sum:         {0}" -f ([System.BitConverter]::ToString($bs).ToLower() -replace "-",""))
}

function Download-Exe($url,$exename){
    if( Test-Path $exename ){
        Write-Verbose -Message ("Found {0}" -f $exename)
        return
    }
    Write-Verbose -Message ("{0} not found." -f $exename)
    Write-Verbose -Message ("$ $GO get " + $url)
    & $GO get $url
    $workdir = (Join-Path (Join-Path (Get-Go1stPath) "src") $url)
    $cwd = (Get-Location)
    Set-Location $workdir
    Write-Verbose -Message ("$ $GO build {0} on {1}" -f $exename,$workdir)
    & $GO build
    Do-Copy $exename $cwd
    Set-Location $cwd
}

function Newer-Than($source,$target){
    if( -not $target ){
        Write-Warning ('Newer-Than: $target is null')
        if( $source ){
            Write-Verbose ('Newer-Than: $source={0}' -f $source)
        }else{
            Write-Warning 'Newer-Than: $source is null'
        }
        return
    }
    if( -not (Test-Path $target) ){
        Write-Verbose ("{0} not found." -f $target)
        return $true
    }
    $stamp = (Get-ItemProperty $target).LastWriteTime

    $prop = (Get-ItemProperty $source)
    if( $prop.Mode -like 'd*' ){
        Get-ChildItem $source -Recurse | %{
            if( $_.LastWriteTime -gt $stamp ){
                Write-Verbose ("{0} is newer than {1}" -f $_.FullName,$target)
                return $true
            }
        }
        return $false
    }else{
        return $prop.LastWriteTime -gt $stamp
    }
}


function Build($version,$tags) {
    Write-Verbose "Build as version='$version' tags='$tags'"

    Go-Generate
    if( -not (Go-Fmt) ){
        return
    }
    $saveGOARCH = $env:GOARCH
    $env:GOARCH = (Get-GoArch)

    Make-Dir $CMD
    $binDir = (Join-Path $CMD $env:GOARCH)
    Make-Dir $binDir
    $target = (Join-Path $binDir "nyagos.exe")

    Make-SysO $version

    $ldflags = (git log -1 --date=short --pretty=format:"-X main.stamp=%ad -X main.commit=%H")
    Write-Verbose "$ $GO build -o '$target'"
    & $GO build "-o" $target -ldflags "$ldflags -X main.version=$version" $tags
    if( $LastExitCode -eq 0 ){
        Do-Copy $target ".\nyagos.exe"
    }
    $env:GOARCH = $saveGOARCH
}

function Make-CSource($xml){
    $const = $xml.make.const
    $package = $const.package

    foreach( $h in $const.include ){
        Write-Output "#include $h"
    }
    Write-Output '#define X(x) #x'
    Write-Output ''
    Write-Output 'int main()'
    Write-Output '{'
    Write-Output ('    printf("package {0}\n\n");' -f $package )

    foreach( $p in $const.li ){
        $name1 = $p.nm
        $type  = $p.type
        $fmt   = $p.fmt
        if( $fmt ){
            Write-Output ('     printf("const " X({0}) "={1}\n",{0});' `
                -f $name1,$fmt)
        }elseif( $type ){
            Write-Output ('     printf("const " X({0}) "={1}(%d)\n",{0});' `
                -f $name1,$type)
        }else{
            Write-Output ('     printf("const " X({0}) "=%d\n",{0});' `
                -f $name1)
        }
    }
    Write-Output '    return 0;'
    Write-Output '}'
}

function Make-ConstGo($makexml){
    Make-CSource $makexml |
        Out-File "makeconst.c" -Encoding Default

    Write-Verbose -Message '$ gcc makeconst.c'
    gcc "makeconst.c"

    if( -not (Test-Path "a.exe") ){
        Write-Error -Message "a.exe not found"
        return
    }
    Write-Verbose -Message '$ .\a.exe > const.go'
    & ".\a.exe" | Out-File const.go -Encoding default
    Write-Verbose -Message "$ $GO fmt const.go"
    & $GO fmt const.go

    Do-Remove "makeconst.o"
    Do-Remove "makeconst.c"
    Do-Remove "a.exe"
}

function Byte2DWord($a,$b,$c,$d){
    return ($a+256*($b+256*($c+256*$d)))
}

function Get-Architecture($bin){
    $addr = (Byte2DWord $bin[60] $bin[61] $bin[62] $bin[63])
    if( $bin[$addr] -eq 0x50 -and $bin[$addr+1] -eq 0x45 ){
        if( $bin[$addr+4] -eq 0x4C -and $bin[$addr+5 ] -eq 0x01 ){
            return 32
        }
        if( $bin[$addr+4] -eq 0x64 -and $bin[$addr+5] -eq 0x86 ){
            return 64
        }
    }
    return $null
}

function Download-File($url){
    $fname = (Split-Path $url -Leaf)
    if( $fname -like "download*" ){
        $fname = (Split-Path (Split-Path $url -Parent) -Leaf)
    }
    $client = New-Object System.Net.WebClient
    Write-Verbose "$ wget '$url' -> '$fname'"
    $client.DownloadFile($url,$fname)
    return $fname
}

function Get-Lua($url,$arch){
    $zip = (Download-File $url)
    unzip -o $zip include\*
    $folder = (Join-Path $CMD $arch)
    Make-Dir $CMD
    Make-Dir $folder
    unzip -o $zip lua53.dll -d $folder
    Do-Copy (Join-Path $folder lua53.dll) .
}

function Make-Package($arch){
    $zipname = ("nyagos-{0}.zip" -f (& cmd\$arch\nyagos.exe --show-version-only))
    Write-Verbose "$ zip -9 $zipname ...."
    if( Test-Path $zipname ){
        Do-Remove $zipname
    }
    zip -9j $zipname `
        "cmd\$arch\nyagos.exe" `
        "cmd\$arch\lua53.dll" `
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
        Build (git describe --tags) ""
    }
    "nolua" {
        Build "" "-tags=nolua"
    }
    "386"{
        $private:save = $env:GOARCH
        $env:GOARCH = "386"
        Build (git describe --tags) ""
        $env:GOARCH = $save
    }
    "debug" {
        $private:save = $env:GOARCH
        if( $args[1] ){
            $env:GOARCH = $args[1]
        }
        Build "" "-tags=debug"
        $env:GOARCH = $save
    }
    "release" {
        $private:save = $env:GOARCH
        if( $args[1] ){
            $env:GOARCH = $args[1]
        }
        Build (Get-Content Misc\version.txt) ""
        $env:GOARCH = $save
    }
    "status" {
        Show-Version ".\nyagos.exe"
    }
    "vet" {
        ForEach-GoDir | ForEach-Object{
            pushd $_
            Write-Verbose "$ $GO vet on $_"
            & $GO vet
            popd
        }
    }
    "clean" {
        foreach( $p in @(`
            (Join-Path $CMD "amd64\nyagos.exe"),`
            (Join-Path $CMD "386\nyagos.exe"),`
            "nyagos.exe",`
            "nyagos.syso",`
            "version.now",`
            "goversioninfo.exe") )
        {
            Do-Remove $p
        }
        Get-ChildItem "." -Recurse |
        Where-Object { $_.Name -eq "make.xml" } |
        ForEach-Object {
            $dir = (Split-Path $_.FullName -Parent)
            $xml = [xml](Get-Content $_.FullName)
            foreach($li in $xml.make.generate.li){
                if( -not $li ){ continue }
                foreach($target in $xml.make.generate.li.target){
                    if( -not $target ){ continue }
                    $path = (Join-Path $dir $target)
                    if( Test-Path $path ){
                        Do-Remove $path
                    }
                }
            }
        }

        ForEach-GoDir | %{
            Write-Verbose "$ $GO clean on $_"
            pushd $_
            & $GO clean
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
        ?{ $_.Name -eq "make.xml" } |
        %{
            $makexml = [xml](Get-Content $_.FullName)
            $private:cwd = (Split-Path $_.FullName -Parent)
            pushd $cwd
            Write-Verbose "$ chdir $cwd"
            $private:pkg = (Split-Path $cwd -Leaf)
            Write-Verbose "for package $pkg"
            Make-ConstGo $makexml
            popd
        }
    }
    "package" {
        $goarch = if( $args[1] ){ $args[1] }else{ (Get-GoArch) }
        Make-Package $goarch
    }
    "install" {
        $installDir = $args[1]
        if( $installDir -eq $null -or $installDir -eq "" ){
            $installDir = (
                Select-String 'INSTALLDIR=([^\)"]+)' Misc\version.cmd |
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
            Out-File "Misc\version.cmd" -Encoding Default

        robocopy nyagos.d (Join-Path $installDir "nyagos.d") /E
        Write-Verbose ("ERRORLEVEL=" + $LastExitCode)
        if( $LastExitCode -lt 8 ){
            Remove-Item Variable:LastExitCode
        }
        if( -not (Test-Path (Join-Path $installDir "lua53.dll") ) ){
            Do-Copy lua53.dll $installDir
        }
        Ask-Copy "_nyagos" $installDir
        try{
            Do-Copy nyagos.exe $installDir
        }catch{
            taskkill /F /im nyagos.exe
            try{
                Do-Copy nyagos.exe $installDir
            }catch{
                Write-Host "Could not update installed nyagos.exe"
                Write-Host "Some processes holds nyagos.exe now"
            }
            # [void]([System.Windows.Forms.MessageBox]::Show("Done"))
            timeout /T 3
        }
    }
    "generate" {
        Go-Generate
    }
    "get" {
        Go-Generate
        Get-Imports | ForEach-Object{ Write-Output $_ ; & $GO get -u $_ }
    }
    "fmt" {
        Go-Fmt | Out-Null
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
                        $private:output = "$one,$other"
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
    "get-lua" {
        $goarch = if( $args[1] ){ $args[1] } else { (Get-GoArch) }
        Get-Lua ($LUAURL[ $goarch ]) $goarch
    }
    "help" {
        Write-Output @'
make                     build as snapshot
make debug   [386|amd64] build as debug version     (tagged as `debug`)
make release [386|amd64] build as release version
make status              show version information 
make vet                 do `go vet` on each folder
make clean               remove all work files
make sweep               remove *.bak and *.~
make const               make `const.go`. gcc is required
make package [386|amd64] make `nyagos-(VERSION)-(ARCH).zip`
make install [FOLDER]    copy executables to FOLDER or last folder
make generate            execute `go generate` on the folder it required
make fmt                 `go fmt`
make check-case [FILE]
make get-lua [386|amd64] download Lua 5.3 for current architecture
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
