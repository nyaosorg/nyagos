local tmpfn = nyagos.env.temp .. "\\tmp.txt"

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
for line in fd:lines() do
    if line == "HOGEHOGE" then
        print("OK: ",line)
    else
        print("NG: ",line)
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
