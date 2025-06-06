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

    -- The function `use` is now deprecated.
end

function use(name)
    require(string.gsub(name,"%.lua$",""))
end
