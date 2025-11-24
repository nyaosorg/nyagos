share.jj={
    ["abandon"]={},
    ["absorb"]={},
    ["bisect"]={
        ["run"]={},
    },
    ["bookmark"]={
        ["create"]={},
        ["delete"]={},
        ["forget"]={},
        ["list"]={},
        ["move"]={},
        ["rename"]={},
        ["set"]={},
        ["track"]={},
        ["untrack"]={},
    },
    ["commit"]={},
    ["config"]={
        ["edit"]={},
        ["get"]={},
        ["list"]={},
        ["path"]={},
        ["set"]={},
        ["unset"]={},
    },
    ["describe"]={},
    ["diff"]={},
    ["diffedit"]={},
    ["duplicate"]={},
    ["edit"]={},
    ["evolog"]={},
    ["file"]={
        ["annotate"]={},
        ["chmod"]={},
        ["list"]={},
        ["show"]={},
        ["track"]={},
        ["untrack"]={},
    },
    ["fix"]={},
    ["gerrit"]={
        ["upload"]={},
    },
    ["git"]={
        ["clone"]={},
        ["colocation"]={},
        ["export"]={},
        ["fetch"]={},
        ["import"]={},
        ["init"]={},
        ["push"]={},
        ["remote"]={},
        ["root"]={},
    },
    ["help"]={},
    ["interdiff"]={},
    ["log"]={},
    ["metaedit"]={},
    ["new"]={},
    ["next"]={},
    ["operation"]={
        ["abandon"]={},
        ["diff"]={},
        ["log"]={},
        ["restore"]={},
        ["revert"]={},
        ["show"]={},
    },
    ["parallelize"]={},
    ["prev"]={},
    ["rebase"]={},
    ["redo"]={},
    ["resolve"]={},
    ["restore"]={},
    ["revert"]={},
    ["root"]={},
    ["show"]={},
    ["sign"]={},
    ["simplify-parents"]={},
    ["sparse"]={
        ["edit"]={},
        ["list"]={},
        ["reset"]={},
        ["set"]={},
    },
    ["split"]={},
    ["squash"]={},
    ["status"]={},
    ["tag"]={
        ["delete"]={},
        ["list"]={},
        ["set"]={},
    },
    ["undo"]={},
    ["unsign"]={},
    ["util"]={
        ["completion"]={},
        ["config-schema"]={},
        ["exec"]={},
        ["gc"]={},
        ["install-man-pages"]={},
        ["markdown-help"]={},
    },
    ["version"]={},
    ["workspace"]={
        ["add"]={},
        ["forget"]={},
        ["list"]={},
        ["rename"]={},
        ["root"]={},
        ["update-stale"]={},
    },
}
nyagos.complete_for["jj"] = function(args)
    if not string.match(args[#args],"^[-a-z]+") then
        return nil
    end

    local j = share.jj
    local last = nil
    while true do
        repeat
            table.remove(args,1)
            if #args <= 0 then
                return last
            end
            last = args[1]
        until string.sub(last,1,1) ~= "-"

        local nextj = j[ last ]
        if not nextj then
            local result = {}
            for key,val in pairs(j) do
                result[#result+1] = key
            end
            if next(result) then
                return result
            else
                return nil
            end
        end
        j = nextj
    end
end
