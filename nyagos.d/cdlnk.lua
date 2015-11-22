if not share.ole then
    local status
    status, share.ole = pcall(require,"nyole")
    if not status then
        share.ole = nil
    end
end

if share.ole then
    share._shortcut_getfolder = function(arg1)
        local fsObj = share.ole.create_object_utf8("Scripting.FileSystemObject")
        local wshObj = share.ole.create_object_utf8("WScript.Shell")

        local shortcut1 = wshObj:CreateShortcut(arg1)
        if not shortcut1 then
            return arg1
        end
        local path = shortcut1.TargetPath
        if fsObj:FolderExists(path) then
            return path
        end
        path = shortcut1.WorkingDirectory
        if fsObj:FolderExists(path) then
            return path
        end
        path = fsObj:GetParentFolderName(shortcut1.TargetPath)
        if fsObj:FolderExists(path) then
            return path
        end
        return arg1
    end

    share.org_cdlnk_alias_cd = nyagos.alias.cd
    nyagos.alias.cd=function(args)
        for i=1,#args do
            local arg1 = args[i]
            if string.match(arg1,"%.[lL][nN][kK]$") then
                arg1 = share._shortcut_getfolder(arg1)
            end
            args[i] = arg1
        end
        if share.org_cdlnk_alias_cd then
            return share.org_cdlnk_alias_cd(args)
        else
            args[0] = "__cd__"
            return nyagos.exec(args)
        end
    end
end
