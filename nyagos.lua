--------------------------------------------------------------------------
-- DO NOT EDIT THIS. PLEASE EDIT ~\.nyagos OR ADD SCRIPT INTO nyagos.d\ --
--------------------------------------------------------------------------

if nyagos == nil then
    print("This is the startup script for NYAGOS")
    print("Do not run this with lua.exe")
    os.exit(0)
end

print(string.format("Nihongo Yet Another GOing Shell %s-%s Powered by %s & %s",
      (nyagos.version or "v"..nyagos.stamp),
      nyagos.goarch,
      nyagos.goversion,_VERSION))

print("Copyright (c) 2014-2016 HAYAMA_Kaoru and NYAOS.ORG")

local function include(fname)
    local chank,err=nyagos.loadfile(fname)
    if err then
        print(err)
    elseif chank then
        local ok,err=pcall(chank)
        if not ok then
            print(fname .. ": " ..err)
        end
    else
        print(fname .. ":fail to load")
    end
end

function use(name)
    local catalog_d = string.gsub(nyagos.exe,"%.[eE][xX][eE]$",".d\\catalog")
    name = string.gsub(name,"%.lua$","") .. ".lua"
    local script = nyagos.pathjoin(catalog_d,name)
    include(script)
end

local addons=nyagos.glob((nyagos.exe:gsub("%.[eE][xX][eE]$",".d\\*.lua")))
for i=1,#addons do
    include(addons[i])
end

local home = nyagos.getenv("HOME") or nyagos.getenv("USERPROFILE")
if home then 
    local dotfile = nyagos.pathjoin(home,'.nyagos')
    local fd=nyagos.open(dotfile)
    if fd then
        fd:close()
        include(dotfile)
    else
        dotfile = string.gsub(nyagos.exe,"[^\\\\\\/]+$",".nyagos")
        fd = nyagos.open(dotfile)
        if fd then
            fd:close()
            include(dotfile)
        end
    end
end
