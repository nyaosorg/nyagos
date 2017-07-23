Set-PSDebug -strict

Function Get-GoArch{
    if( Test-Path "goarch.txt" ){
        (Get-Content "goarch.txt")
    }else{
        ( go version | %{ $_.Split()[-1].Split("/")[-1] })
    }
}

Function Get-Version{
    Get-Content Misc\version.txt
}

Function Get-Commit{
    git describe --tags
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

function Get-Modified {
    Get-ChildItem . -Recurse -Name |
    ?{ $_ -like "*.go" } |
    ?{ 
        $private:attr = [System.IO.File]::GetAttributes($_)
        ($attr -band [System.IO.FileAttributes]::Archive) -eq
            [System.IO.FileAttributes]::Archive
    }
}

function Set-NoModified($path) {
    $private:attr = [System.IO.File]::GetAttributes($path)
    [System.IO.File]::SetAttributes($path,
        $attr -band -bnot [System.IO.FileAttributes]::Archive)
}

function Go-Fmt{
    Get-Modified | %{ go fmt -n $_ ; Set-NoModified($_) }
}

function Get-Go1stPath {
    if( $env:gopath -ne "" ){
        $gopath = $env:gopath
    }else{
        $gopath = Join-Path($env:userprofile,"go")
    }
    $gopath.Split(";")[0]
}

function Make-JSON($version){
    $version = [string]$version
    $json = @{
        FixedFileInfo = @{
            FileFlagsMask="3f";
            FileFlags="00";
            FileOS="040004";
            FileType="01";
            FileSubType="00";
        };
        StringFileInfo = @{
            Comments= "";
            CompanyName= "NYAOS.ORG";
            FileDescription= "Extended Commandline Shell";
            InternalName= "";
            LegalCopyright= "Copyright (C) 2014-2017 HAYAMA_Kaoru";
            LegalTrademarks= "";
            OriginalFilename= "NYAGOS.EXE";
            PrivateBuild= "";
            ProductName= "Nihongo Yet Another GOing Shell";
            SpecialBuild= "";
        };
        VarFileInfo = @{
            Translation=@{
                LangID= "0411";
                CharsetID= "04E4";
            }
        }
    }
    if( $version -like "*.*.*_*" ){
        $v = $version -split "[\._]"
    }else{
        $v = @(0,0,0,0)
    }
    $private:versionSet = @{
        Major = $v[0];
        Minor = $v[1];
        Patch = $v[2];
        Build = $v[3];
    }
    $json["FixedFileInfo"]["FileVersion"] = $versionSet
    $json["FixedFileInfo"]["ProductVersion"] = $versionSet
    $json["StringFileInfo"]["FileVersion"] = $version
    $json["StringFileInfo"]["ProductVersion"] = $version

    $json | ConvertTo-Json
}

function Go-VersionInfo($version) {
    $cwd = Get-Location
    if( -not (Test-Path("goversioninfo.exe")) ){
        go get "github.com/josephspurrier/goversioninfo"
        $gopath1 = Get-Go1stPath
        pushd (Join-Path $gopath1 "src\github.com\josephspurrier\goversioninfo\cmd\goversioninfo")
        go build
        copy goversioninfo.exe $cwd
        popd
    }
    Make-JSON( $version ) | Set-Content -encoding UTF8 Misc\version.json
    goversioninfo.exe -icon mains\nyagos.ico -o nyagos.syso Misc\version.json
}

function Show-Version($fname) {
    if( Test-Path $fname ){
        Write-Host $fname
        $v = [System.Diagnostics.FileVersionInfo]::GetVersionInfo($fname)
        if( $v ){
            Write-Host(("  FileVersion:    `"{0}`" ({1},{2},{3},{4})" -f
                $v.FileVersion,
                $v.FileMajorPart,
                $v.FileMinorPart,
                $v.FileBuildPart,
                $v.FilePrivatePart))
            Write-Host(("  ProductVersion: `"{0}`" ({1},{2},{3},{4})" -f
                $v.ProductVersion,
                $v.ProductMajorPart,
                $v.ProductMinorPart,
                $v.ProductBuildPart,
                $v.ProductPrivatePart))
        }
        $data = [System.IO.File]::ReadAllBytes($fname)
        $md5 = New-Object System.Security.Cryptography.MD5CryptoServiceProvider
        $bs = $md5.ComputeHash($data)
        Write-Host( ("  md5sum: {0}" -f [System.BitConverter]::ToString($bs).ToLower().Replace("-","")))
    }else{
        Write-Host(("{0}: not found" -f $fname))
    }
}

function Make-Bindata {
    # Get-ChildItem .\nyagos.d -Recurse |
}

function Do-Build {

}


switch( $args[0] ){
    "status" {
        Show-Version(".\nyagos.exe")
    }
    "" {
        Write-Host("goarch=" + (Get-GoArch))
        Write-Host("version=" + (Get-Version))
        Write-Host("commit=" + (Get-Commit))
        EachGoDir | %{ Write-Host $_ } 
        Get-Imports | %{ Write-Host $_ }
    }
    "get" {
        Get-Imports | %{ Write-Host $_ ; go get -u $_ }
    }
    "fmt" {
        Go-Fmt
    }
    "t" {
        Go-VersionInfo
        #Make-JSON("1.2.3_4")
    }
}
