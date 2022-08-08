
local val,err = nyagos.commonprefix{ "ABC123" , "ABCCCC" }
if not val or val ~= "ABC" then
    os.exit(1)
end

local val,err = nyagos.commonprefix()
if val then
    os.exit(1)
end

local val,err = nyagos.commonprefix(1)
if val then
    os.exit(1)
end
