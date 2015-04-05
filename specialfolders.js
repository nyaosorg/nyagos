var objShell = new ActiveXObject("WScript.Shell");
for(var i=0 ; i < WScript.Arguments.length ; i++ ){
    WScript.Echo( objShell.SpecialFolders(WScript.Arguments.Item(i) ) )
}
