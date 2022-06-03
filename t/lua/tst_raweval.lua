-- lua_f t_raweval.lua --
local x = nyagos.raweval('cmd','/c','echo AHAHA')
if x ~= "AHAHA\r\n" then
    os.exit(1)
end
