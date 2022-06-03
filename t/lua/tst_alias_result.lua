--- test for alias function's return value ---
--- Do `lua_f 20180504.lua`

nyagos.alias.tst20180504 = function(args)
    return { "echo","ahaha" }
end

local result = nyagos.eval("tst20180504")
if result ~= "ahaha" then
    print("NG:",result)
end

nyagos.alias.tst20180504b = function(args)
    return "echo ihihi"
end

result = nyagos.eval("tst20180504b")
if result ~= "ihihi" then
    print("NG:",result)
end

nyagos.alias.tst20180504c = function(args)
    return 3
end

nyagos.exec("tst20180504c")
result = nyagos.env.errorlevel 
if result ~= "3" then
    print("NG",result)
end

