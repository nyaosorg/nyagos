Set-PSDebug -strict

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
            Comments="";
            CompanyName="NYAOS.ORG";
            FileDescription="Extended Commandline Shell";
            InternalName="";
            LegalCopyright="Copyright (C) 2014-2017 HAYAMA_Kaoru";
            LegalTrademarks="";
            OriginalFilename="NYAGOS.EXE";
            PrivateBuild="";
            ProductName="Nihongo Yet Another GOing Shell";
            SpecialBuild="";
        };
        VarFileInfo = @{
            Translation=@{
                LangID="0411";
                CharsetID="04E4";
            }
        }
    }
    $v = $version -split "[\._]"
    if( $v.Length -ne 4 ){
        $v = @(0,0,0,0)
    }
    $private:versionSet = @{
        Major = [int]$v[0];
        Minor = [int]$v[1];
        Patch = [int]$v[2];
        Build = [int]$v[3];
    }
    $json["FixedFileInfo"]["FileVersion"] = $versionSet
    $json["FixedFileInfo"]["ProductVersion"] = $versionSet
    $json["StringFileInfo"]["FileVersion"] = $version
    $json["StringFileInfo"]["ProductVersion"] = $version

    $json | ConvertTo-Json
}

Make-JSON( [string]$env:version )
