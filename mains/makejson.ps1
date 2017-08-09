Set-PSDebug -strict

$text = @"
{
	"FixedFileInfo":
	{
		"FileVersion": {
			"Major": %VER1%,
			"Minor": %VER2%,
			"Patch": %VER3%,
			"Build": %VER4%
		},
		"ProductVersion": {
			"Major": %VER1%,
			"Minor": %VER2%,
			"Patch": %VER3%,
			"Build": %VER4%
		},
		"FileFlagsMask": "3f",
		"FileFlags ": "00",
		"FileOS": "040004",
		"FileType": "01",
		"FileSubType": "00"
	},
	"StringFileInfo":
	{
		"Comments": "",
		"CompanyName": "NYAOS.ORG",
		"FileDescription": "Extended Commandline Shell",
		"FileVersion": "%VERSION%",
		"InternalName": "",
		"LegalCopyright": "Copyright (C) 2014-2017 HAYAMA_Kaoru",
		"LegalTrademarks": "",
		"OriginalFilename": "NYAGOS.EXE",
		"PrivateBuild": "",
		"ProductName": "Nihongo Yet Another GOing Shell",
		"ProductVersion": "%VERSION%",
		"SpecialBuild": ""
	},
	"VarFileInfo":
	{
		"Translation": {
			"LangID": "0411",
			"CharsetID": "04E4"
		}
	}
}
"@

$version = $env:version 
if( $version -like "*.*.*_*" ){
    $v = $version -split "[\._]"

    $text = $text -replace "%VER1%",$v[0]
    $text = $text -replace "%VER2%",$v[1]
    $text = $text -replace "%VER3%",$v[2]
    $text = $text -replace "%VER4%",$v[3]
    $text = $text -replace "%VERSION%",$version
}else{
    $text = $text -replace "%VER[1-4]%","0"
    $text = $text -replace "%VERSION%","0.0.0_0"
}
Write-Host($text)
