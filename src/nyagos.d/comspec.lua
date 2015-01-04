for _,name in pairs{
    "assoc",
    "attrib",
    "dir",
    "mklink",
    "ren",
    "rename",
    "start",
    "type" ,
} do
    nyagos.alias(name,"%COMSPEC% /c "..name.." $*")
end
