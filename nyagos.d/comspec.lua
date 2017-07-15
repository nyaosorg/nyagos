if not nyagos then
    print("This is a script for nyagos not lua.exe")
    os.exit()
end

for _,name in pairs{
    "assoc",
    "dir",
    "mklink",
    "ren",
    "rename",
} do
    nyagos.alias[name] = "%COMSPEC% /c "..name.." $*"
end

nyagos.alias.grep = "findstr.exe"
