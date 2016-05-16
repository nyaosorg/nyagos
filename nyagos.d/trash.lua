nyagos.alias.trash = function(args)
    if #args <= 0 then
        nyagos.writerr("Move files or directories to Windows Trashbox\n")
        nyagos.writerr("Usage: trash file(s)...\n")
        return
    end
    local fsObj = nyagos.create_object("Scripting.FileSystemObject")
    local shellApp = nyagos.create_object("Shell.Application")
    local trashBox = shellApp:_call("NameSpace",math.tointeger(10))
    args = nyagos.glob(table.unpack(args))
    for i=1,#args do
        if fsObj:_call("FileExists",args[i]) or fsObj:_call("FolderExists",args[i]) then
            trashBox:_call("MoveHere",fsObj:_call("GetAbsolutePathName",args[i]))
        else
            nyagos.writerr(args[i]..": such a file or directory not found.\n")
        end
    end
end
