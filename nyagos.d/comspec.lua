for _,name in pairs{
    "assoc",
    "dir",
    "mklink",
    "ren",
    "rename",
} do
    nyagos.alias[name] = "%COMSPEC% /c "..name.." $*"
end
