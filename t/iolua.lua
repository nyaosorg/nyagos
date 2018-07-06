local tmpfn = nyagos.env.temp .. "\\tmp.txt"

local fd = assert(io.open(tmpfn,"w"))
assert(fd:write('HOGEHOGE\n'))
fd:close()

local ok=true
for line in io.lines(tmpfn) do
    if line ~= 'HOGEHOGE' then
        ok = false
    end
end
if ok then
    print("OK: write-open,write-close,io-lines")
else
    print("NG: write-open,write-close,io-lines")
end

local fd,msg = io.open(":::","w")
if fd then
    print("NG: invalid-open")
else
    print("OK: invalid-open:",msg)
end
