if( WScript.Arguments.length == 1  ){
    var wshShell=new ActiveXObject("WScript.Shell");
    var lnkSrc = wshShell.CreateShortcut(WScript.Arguments.Unnamed(0));
    WScript.Echo( "    " + lnkSrc + "\n<-- " + lnkSrc.TargetPath );
    WScript.Quit(1);
}
if( WScript.Arguments.length < 2 ){
    WScript.Echo(
        "Usage: cscript lnk.js FILENAME SHORTCUT {Option=Value}... make shortcut\n" +
        "       cscript lnk.js SHORTCUT          ... print shortcut-target")
    WScript.Quit(1);
}

var fsObj = new ActiveXObject("Scripting.FileSystemObject");
var src = fsObj.GetAbsolutePathName(WScript.Arguments.Item(0));
var dst = fsObj.GetAbsolutePathName(WScript.Arguments.Item(1));

if( fsObj.FolderExists(dst) ){
    dst = fsObj.BuildPath(dst,fsObj.GetFileName(src));
}
if( dst.length >= 4 && dst.substring(dst.length-4) != ".lnk" ){
    dst += ".lnk";
}
var wshShell=new ActiveXObject("WScript.Shell");
var shortcut1=wshShell.CreateShortcut(dst);

if( shortcut1 == null ){
    WScript.Echo("Fail to create ShortCut Object");
    WScript.Quit(1);
}
shortcut1.TargetPath=src;

if( WScript.Arguments.length >= 3 ){
    for(var i=2 ; i < WScript.Arguments.length ; i++ ){
        var equation = WScript.Arguments.Item(i);
        var pos = equation.indexOf("=",0);
        if( pos >= 0 ){
            equation="shortcut1." + equation.substring(0,pos) + "=\"" +
                    equation.substring(pos+1).replace(/\\/g,"\\\\") + "\"";
            WScript.Echo(equation);
            eval(equation);
        }else{
            WScript.Echo("Equal(=) not found: " + equation)
        }
    }
}
shortcut1.Save()
WScript.Echo("    " + src + "\n--> " + dst);
WScript.Quit(0)
