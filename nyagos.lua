--------------------------------------------------------------------------
-- DO NOT EDIT THIS. PLEASE EDIT ~\.nyagos OR ADD SCRIPT INTO nyagos.d\ --
--------------------------------------------------------------------------

print("Nihongo Yet Another GOing Shell")
print("Build at ".. nyagos.stamp .. " with commit "..nyagos.commit)
print("Copyright (c) 2014 HAYAMA_Kaoru and NYAOS.ORG")

for _,fname in ipairs(nyagos.glob("nyagos.d\\*.lua")) do
    local chank,err=assert(loadfile(fname))
    if err then
        print(err)
    else
        chank()
    end
end

local home = nyagos.getenv("HOME") or nyagos.getenv("USERPROFILE")
if home then
    local rcfname = home .. '\\.nyagos'
    local chank,err=assert(loadfile(rcfname))
    if err then
        print(err)
    else
        chank()
    end
end
