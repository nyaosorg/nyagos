nyagos.on_command_not_found = function(args)
    nyagos.writerr(args[0]..": コマンドではない。\n")
    return true
end

if not nyagos.ole then
    local status
    status, nyagos.ole = pcall(require,"nyole")
    if not status then
        nyagos.ole = nil
    end
end

if nyagos.ole then
    local fsObj = nyagos.ole.create_object_utf8("Scripting.FileSystemObject")
    local cd = nyagos.alias.cd
    nyagos.alias.cd = function(args)
        local success=true
        for i=1,#args do
            local dir=args[i]
            if dir:sub(1,1) ~= "-" and not dir:match("%.[lL][nN][kK]$") then
                if not fsObj:FolderExists(dir) then
                    nyagos.writerr(dir..": ディレクトリではない。\n")
                    return
                end
            end
        end
        if cd then
            cd(args)
        else
            args[0] = "__cd__"
            nyagos.exec(args)
        end
    end
end
-- vim:set fenc=utf8 --
