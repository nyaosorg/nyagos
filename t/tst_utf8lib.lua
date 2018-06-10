
--- Compare
---   lua_f tst_utf8lib.lua
--  and
---   lua.exe tst_utf8lib.lua | nkf32

for i,c in utf8.codes("あいうえお") do
    print(i,c,type(c))
end

for c in string.gmatch("あいうえお",utf8.charpattern) do
    print(c)
end
