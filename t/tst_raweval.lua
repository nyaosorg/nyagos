-- lua_f t_raweval.lua --
x = nyagos.raweval('cmd','/c','echo AHAHA')
if x == "AHAHA\r\n" then
    print("OK")
else
    print("NG: x==["..x.."]")
end
