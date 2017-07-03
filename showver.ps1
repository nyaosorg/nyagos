foreach($fname in $args ){
    if( Test-Path $fname ){
        $v = [System.Diagnostics.FileVersionInfo]::GetVersionInfo($fname)
        if( $v ){
            Write-Host(("FileVersion:    {0} ({1},{2},{3},{4})" -f
                $v.FileVersion,
                $v.FileMajorPart,
                $v.FileMinorPart,
                $v.FileBuildPart,
                $v.FilePrivatePart))
            Write-Host(("ProductVersion: {0} ({1},{2},{3},{4})" -f
                $v.ProductVersion,
                $v.ProductMajorPart,
                $v.ProductMinorPart,
                $v.ProductBuildPart,
                $v.ProductPrivatePart))
        }
    }else{
        Write-Host(("{0}: not found" -f $fname))
    }
}
