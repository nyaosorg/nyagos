if not nyagos then
    print("This is a script for nyagos not lua.exe")
    os.exit()
end

if nyagos.goos == "windows" then
    local comspec = nyagos.env.comspec
    if not comspec or comspec == "" then
        comspec = "CMD.EXE"
    end
    for _,name in pairs{
        "assoc",
        "dir",
        "mklink",
        "ren",
        "rename",
    } do
        nyagos.alias[name] = comspec .. " /c "..name.." $*"
    end

    local greppath=nyagos.which("grep")
    if not greppath and not nyagos.alias.grep then
        nyagos.alias.grep = "findstr.exe"
    end
end
