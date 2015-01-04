function findValue(name,dir)
    local cmdline = string.format('findstr "\\<%s\\>" %s\\*.h',name,dir)
    local fd = io.popen(cmdline,"r")
    local value=nil
    for line in fd:lines() do
        for token in string.gmatch(line,"%S+") do
            value=token
        end
        if value then
            return value
        end
    end
    return ""
end
function output(keyword,includedir,out)
    for _,keyword in ipairs(keyword) do
        local value = findValue(keyword,includedir)
        if value then
            value = string.gsub(value,"(DWORD)","uint32")
            print(string.format("const %s=%s",keyword,value))
        end
    end
end
if #arg < 3 then
    print "Usage: lua makeconst.lua KEYWORD(s)... INCLUDEDIR"
    os.exit()
end

local includedir=arg[#arg]
local package=arg[1]
table.remove(arg,1)
table.remove(arg,#arg)

print("package "..package)
output(arg,includedir)
