var fsObj=new ActiveXObject("Scripting.FileSystemObject");
for(var i=0 ; i < WScript.Arguments.length ; i++ ){
    var path=fsObj.GetAbsolutePathName(WScript.Arguments.Item(i));
    if( ! fsObj.FileExists(path) ){
        continue
    }
    var version=fsObj.GetFileVersion(path);
    if( WScript.Arguments.length == 1 ){
        WScript.Echo(version)
    }else{
        WScript.Echo(path+" "+version)
    }
}
