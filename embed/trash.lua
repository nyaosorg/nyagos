if not nyagos then
    print("This is a script for nyagos not lua.exe")
    os.exit()
end

nyagos.alias.trash = function(args)
    if #args <= 0 then
        nyagos.writerr("Move files or directories to Windows Trashbox\n")
        nyagos.writerr("Usage: trash file(s)...\n")
        return
    end
    local fsObj = nyagos.create_object("Scripting.FileSystemObject")
    local shellApp = nyagos.create_object("Shell.Application")
    local trashBox = shellApp:NameSpace(nyagos.to_ole_integer(10))
    args = nyagos.glob((table.unpack or unpack)(args))
    for i=1,#args do
        if fsObj:FileExists(args[i]) or fsObj:FolderExists(args[i]) then
            trashBox:MoveHere(fsObj:GetAbsolutePathName(args[i]))
        else
            nyagos.writerr(args[i]..": such a file or directory not found.\n")
        end
    end
    trashBox:_release()
    shellApp:_release()
    fsObj:_release()
end
