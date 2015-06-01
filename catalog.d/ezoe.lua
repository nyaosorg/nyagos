nyagos.on_command_not_found = function(args)
    nyagos.writerr(args[0]..": コマンドではない。\n")
    return true
end

local cd = nyagos.alias.cd
nyagos.alias.cd = function(args)
    local success=true
    for i=1,#args do
        local dir=args[i]
        if dir:sub(1,1) ~= "-" and not dir:match("%.[lL][nN][kK]$") then
            local nul=nyagos.pathjoin(dir,"nul")
            local fd=io.open(nul)
            if fd then
                fd:close()
            else
                nyagos.writerr(dir..": ディレクトリではない。\n")
                success = false
            end
        end
    end
    if success then
        if cd then
            cd(args)
        else
            args[0] = "__cd__"
            nyagos.exec(args)
        end
    end
end

-- vim:set fenc=utf8 --
