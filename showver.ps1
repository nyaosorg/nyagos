$v = [System.Diagnostics.FileVersionInfo]::GetVersionInfo($Args[0]) 
if( $v ){
    Write-Output(("FileVersion:    {0} ({1},{2},{3},{4})" -f 
        $v.FileVersion,
        $v.FileMajorPart,
        $v.FileMinorPart,
        $v.FileBuildPart,
        $v.FilePrivatePart))
    Write-Output(("ProductVersion: {0} ({1},{2},{3},{4})" -f 
        $v.ProductVersion,
        $v.ProductMajorPart,
        $v.ProductMinorPart,
        $v.ProductBuildPart,
        $v.ProductPrivatePart))
}
