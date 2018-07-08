-- This script runs on both lua.exe and nyagos' lua_f

local tmpfn = os.getenv("TEMP") .. "\\tmp.txt"

local fd = assert(io.open(tmpfn,"w"))
assert(fd:write('HOGEHOGE\n'))
assert(fd:flush())
fd:close()

local ok=true
for line in io.lines(tmpfn) do
    if line ~= 'HOGEHOGE' then
        ok = false
    end
end
if ok then
    print("OK: write-open,write-flush,write-close,io-lines")
else
    print("NG: write-open,write-flush,write-close,io-lines")
end

local fd,err = io.open(":::","w")
if fd then
    print("NG: invalid-open")
else
    print("OK: invalid-open:",err)
end

local fd,err = io.open(tmpfn,"r")
local line = fd:read("*l")
if line == "HOGEHOGE" then
    print("OK: io.read('*l')",line)
else
    print("NG: io.read('*l')",line)
end
fd:seek("set",0)
for line in fd:lines() do
    if line == "HOGEHOGE" then
        print("OK: io.open: ",line)
    else
        print("NG: io.open: ",line)
    end
end
fd:close()

local fd,err = io.popen("cmd.exe /c \"echo AHAHA\"","r")
if fd then
    for line in fd:lines() do
        print("OK>",line)
    end
    fd:close()
else
    print("NG: ",err)
end

local sample=string.gsub(arg[0],"%.lua$",".txt")
print(sample)
local fd,err = io.open(sample,"r")
local line,num,crlf,rest = fd:read("*l","*n",1,"*a")
if io.type(fd) == "file" then
    print "OK: iotype()==\"file\""
else
    print("NG: iotype()==\""..io.type(fd).."\"")
end
fd:close()
if io.type(fd) == "closed file" then
    print "OK: iotype()==\"closed file\""
else
    print("NG: iotype()==\""..io.type(fd).."\"")
end

if io.type("") == nil then
    print "OK: iotype()==nil"
else
    print("NG: iotype()==\""..io.type("").."\"")
end

if line == "ONELINE" then
    print"OK: read('*l')"
else
    print("NG: read('*l'):["..line.."]")
end
if num == 4 then
    print "OK: read('*n')"
else
    print("NG: read('*n')",num)
end
if crlf == "\n" then
    print "OK: read(1)"
else
    print ("NG: read(1):"..crlf)
end
if rest == "AHAHA\nIHIHI\nUFUFU\n" then
    print "OK: read('*a')"
else
    print("NG: read('*a'):["..rest.."]")
end
