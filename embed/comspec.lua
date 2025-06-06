if not nyagos then
    print("This is a script for nyagos not lua.exe")
    os.exit()
end

if nyagos.goos == "windows" then
    for _,name in pairs{
        "assoc",
        "dir",
        "for",
        "ren",
        "rename",
        "date",
        "time"
    } do
        nyagos.alias[name] = "cmdexesc " .. name .. " $*"
    end
end
