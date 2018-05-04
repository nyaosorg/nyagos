--- Do 'lua_f THIS_SCRIPT'

-- for Lua 5.3
--    io.write( (nyagos.utoa('\xE5\x85\x83\x55\x54\x46\x38\xE6\x96\x87\xE5\xAD\x97\xE5\x88\x97\xE3\x82\x92\x53\x4A\x49\x53\xE3\x81\xAB\xE5\xA4\x89\xE6\x8F\x9B\xE3\x81\x99\xE3\x82\x8B'))
-- for Lua 5.1

local hextable = "0123456789ABCDEF"
local dump = {}
for _,val in ipairs{'E5','85','83','55','54','46','38','E6','96','87','E5','AD','97','E5','88','97','E3','82','92','53','4A','49','53','E3','81','AB','E5','A4','89','E6','8F','9B','E3','81','99','E3','82','8B'} do
    dump[ #dump + 1] =
        (string.find(hextable,string.sub(val,1,1))-1) * 16 +
        (string.find(hextable,string.sub(val,2))-1)
end
local utf8 = string.char((table.unpack or unpack)(dump))
local ansi = nyagos.utoa(utf8)
local utf8new = nyagos.atou(ansi)
io.write( utf8new )
