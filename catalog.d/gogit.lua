-- Follow when you miss go and git.

local git_only = {
    commit=true, push=true, pull=true, diff=true, status=true, log=true,
    add=true, rebase=true
}
local go_only = {
    fmt=true, build=true
}

nyagos.alias.go = function(args)
    if #args >= 1 and git_only[args[1]] then
        args[0] = "git" 
    else
        args[0] = "go"
    end
    assert(nyagos.rawexec(args[0],table.unpack(args)))
end
nyagos.alias.git = function(args)
    if #args >=1 and go_only[args[1]] then
        args[0] = "go"
    else
        args[0] = "git"
    end
    assert(nyagos.rawexec(args[0],table.unpack(args)))
end
