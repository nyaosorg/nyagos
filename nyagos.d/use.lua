if not nyagos then
    print("This is a script for nyagos not lua.exe")
    os.exit()
end

do
    local newpath = string.gsub(nyagos.exe,"%.[eE][xX][eE]$",".d\\catalog\\?.lua")
    if package.path and package.path ~= "" then
        package.path = package.path .. ";" .. newpath
    else
        package.path = newpath
    end
end

function use(name)
    local catalog_d = string.gsub(nyagos.exe,"%.[eE][xX][eE]$",".d\\catalog")
    name = string.gsub(name,"%.lua$","") .. ".lua"
    local fname = nyagos.pathjoin(catalog_d,name)
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
