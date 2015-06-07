if not nyagos.ole then
    local status
    status,nyagos.ole = pcall(require,"nyole")
    if not status then
        nyagos.ole = nil
    end
end
if nyagos.ole then
    local fsObj = nyagos.ole.create_object_utf8("Scripting.FileSystemObject")
    local shellApp = nyagos.ole.create_object_utf8("Shell.Application")
    local trashBox = shellApp:NameSpace(math.tointeger(10))
    if trashBox.MoveHere then
        nyagos.alias.trash = function(args)
            if #args <= 0 then
                nyagos.writerr("Move files or directories to Windows Trashbox\n")
                nyagos.writerr("Usage: trash file(s)...\n")
                return
            end
            args = nyagos.glob(table.unpack(args))
            for i=1,#args do
                if fsObj:FileExists(args[i]) or fsObj:FolderExists(args[i]) then
                    trashBox:MoveHere(fsObj:GetAbsolutePathName(args[i]))
                else
                    nyagos.writerr(args[i]..": such a file or directory not found.\n")
                end
            end
        end
    else
        nyagos.writerr("Warning: trash.lua requires nyaole.dll 0.0.0.5 or later\n")
    end
end
