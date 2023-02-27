-- lua_f t_raweval.lua --
local expect
local result

if nyagos.env.OS == "Windows_NT" then
    result = nyagos.raweval('cmd','/c','echo AHAHA')
    expect = "AHAHA\r\n"
else
    result = nyagos.raweval('echo','AHAHA')
    expect = "AHAHA\n"
end

if result ~= expect then
    print(string.format("expect '%s', but '%s'",expect,result))
    os.exit(1)
end
