share.jj={
    ["abandon"]={},
    ["backout"]={},
    ["branch"]={
        ["create"]={},
        ["delete"]={},
        ["forget"]={},
        ["list"]={},
        ["rename"]={},
        ["set"]={},
        ["track"]={},
        ["untrack"]={},
        ["help"]={},
    },
    ["cat"]={},
    ["chmod"]={},
    ["commit"]={},
    ["config"]={
        ["list"]={},
        ["get"]={},
        ["set"]={},
        ["edit"]={},
        ["path"]={},
        ["help"]={},
    },
    ["describe"]={},
    ["diff"]={},
    ["diffedit"]={},
    ["duplicate"]={},
    ["edit"]={},
    ["files"]={},
    ["git"]={
        ["remote"]={},
        ["init"]={},
        ["fetch"]={},
        ["clone"]={},
        ["push"]={},
        ["import"]={},
        ["export"]={},
        ["help"]={},
    },
    ["init"]={},
    ["interdiff"]={},
    ["log"]={},
    ["move"]={},
    ["new"]={},
    ["next"]={},
    ["obslog"]={},
    ["operation"]={
        ["abandon"]={},
        ["log"]={},
        ["undo"]={},
        ["restore"]={},
        ["help"]={},
    },
    ["prev"]={},
    ["rebase"]={},
    ["resolve"]={},
    ["restore"]={},
    ["root"]={},
    ["show"]={},
    ["sparse"]={
        ["list"]={},
        ["set"]={},
        ["help"]={},
    },
    ["split"]={},
    ["squash"]={},
    ["status"]={},
    ["tag"]={
        ["list"]={},
        ["help"]={},
    },
    ["util"]={
        ["completion"]={},
        ["gc"]={},
        ["mangen"]={},
        ["markdown-help"]={},
        ["config-schema"]={},
        ["help"]={},
    },
    ["undo"]={},
    ["unsquash"]={},
    ["untrack"]={},
    ["version"]={},
    ["workspace"]={
        ["add"]={},
        ["forget"]={},
        ["list"]={},
        ["root"]={},
        ["update-stale"]={},
        ["help"]={},
    },
    ["help"]={},
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
