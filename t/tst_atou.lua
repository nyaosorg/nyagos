local hextable = "0123456789ABCDEF"
local dump = {}
for _,val in ipairs{'95','B6','8E','9A','97','F1'} do
    dump[ #dump + 1 ] =
        (string.find(hextable,string.sub(val,1,1))-1)*16 +
        (string.find(hextable,string.sub(val,2))-1)
end
local sjis = 'SHIFTJIS' .. string.char((table.unpack or unpack)(dump))
nyagos.write( nyagos.atou(sjis),"\n" )
