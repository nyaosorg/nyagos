-- This script runs on both lua.exe and nyagos' lua_f

local tmpfn = os.tmpname()

local fd = assert(io.open(tmpfn,"w"))
assert(fd:write('HOGEHOGE\n'))
assert(fd:flush())
fd:close()

local ok=false
for line in io.lines(tmpfn) do
    if line == 'HOGEHOGE' then
        ok = true
    end
end
if not ok then
    print("NG: write-open,write-flush,write-close,io-lines")
    os.exit(1)
end

local fd,err = io.open("./notexistdir/cantopenfile","w")
if fd then
    print("NG: invalid-open")
    fd:close()
    os.exit(1)
end

local fd,err = io.open(tmpfn,"r")
local line = fd:read("*l")
if line ~= "HOGEHOGE" then
    print("NG: io.read('*l')",line)
    os.exit(1)
end
fd:seek("set",0)
for line in fd:lines() do
    if line ~= "HOGEHOGE" then
        print("NG: io.open: ",line)
        os.exit(1)
    end
end
fd:close()

local fd,err = io.popen("echo AHAHA","r")
if fd then
    for line in fd:lines() do
    end
    fd:close()
else
    print("NG: ",err)
    os.exit(1)
end

local sample=string.gsub(arg[0],"%.lua$",".txt")
--print(sample)
local fd,err = io.open(sample,"r")
local line,num,crlf,rest = fd:read("*l","*n",1,"*a")
if io.type(fd) ~= "file" then
    print("NG: iotype()==\""..io.type(fd).."\"")
    os.exit(1)
end

if io.type("") ~= nil then
    print("NG: iotype()==\""..io.type("").."\"")
    os.exit(1)
end

if line ~= "ONELINE" then
    print("NG: read('*l'):[",line,"]")
    os.exit(1)
end
if num ~= 4 then
    print("NG: read('*n')",num)
    os.exit(1)
end
if crlf ~= "\n" then
    print ("NG: read(1):",crlf)
    os.exit(1)
end
if rest ~= "AHAHA\nIHIHI\nUFUFU" then
    print("NG: read('*a'):[",rest,"]")
    os.exit(1)
end

local pos,err = fd:seek('set', 0)
if pos ~= 0 then
  print("NG: seek('set',0)",pos,err)
  os.exit(1)
end
------

local line = fd:read('*l')
if line ~= "ONELINE" then
    print("NG: seek('set',0)", line)
end

local pos,err = fd:seek('set', 9)
if pos ~= 9 then
  print("NG: seek('set',9)==9:",pos)
  os.exit(1)
end

local line = fd:read('*l')
if line ~= '4' then
    print("NG: seek('set',9)==9 read('*l')==4:",line)
    os.exit(1)
end

local pos,err = fd:seek('cur')
if pos ~= 12  then
    print("NG: seek('cur')==12:", pos)
    os.exit(1)
end

local line = fd:read('*l')
if line ~= 'AHAHA' then
    print("NG: seek('cur')==12 read('*l')=='AHAHA'",line)
    os.exit(1)
end

local pos,err = fd:seek('cur',2)
if pos ~= 21 then
  print("NG: seek('cur',2)==21:",pos)
  os.exit(1)
end

local line = fd:read('*l')
if line ~= 'IHI' then
    print("NG: seek('cur',2)==21 read('*l')=='IHI'",line)
    os.exit(1)
end

local pos,err = fd:seek('end')
if pos ~= 31 then
    print("NG: seek('end')==31",pos)
    os.exit(1)
end
local empty = fd:read('*a')
if empty ~= '' then
    print("NG: seek('end')==31 read('*a')==''",empty)
    os.exit(1)
end

local size,err = fd:seek('end')
fd:seek('set',0)
rest,err = fd:read(size + 1)
if rest ~= "ONELINE\n4\nAHAHA\nIHIHI\nUFUFU" then
    print("NG: read(number) number > file size:", rest)
    os.exit(1)
end

fd:seek('set',0)
local all = fd:read('*a')
local before = all:match('.*\n'):gsub('\n','\r\n')
fd:seek('set',#before)
local line = fd:read('*l')
if line ~= 'UFUFU' then
    print("NG: read('*l') on line without line break:",line)
    os.exit(1)
end

fd:seek('end')
local eof = fd:read('*a')
if eof ~= '' then
    print("NG: read('*a') on EOF:",eof)
    os.exit(1)
end

local eof = fd:read('*l')
if eof ~= nil then
    print("NG: read('*l') on EOF:["..eof..']')
    os.exit(1)
end

local eof = fd:read('*n')
if eof then
    print("NG: read('*n') on EOF:["..eof..']')
    os.exit(1)
end

local eof = fd:read(1)
if eof then
    print("NG: read('number') on EOF:["..eof.."]")
    os.exit(1)
end

fd:seek('set',0)
local empty = fd:read(0)
if empty ~= '' then
    print("NG: read(0) on not EOF:",eof)
    os.exit(1)
end
local _, eof = fd:read('*a', 0)
if eof then
    print("NG: read(0) on EOF:",eof)
    os.exit(1)
end

fd:close()
if io.type(fd) ~= "closed file" then
    print("NG: iotype()==\""..io.type(fd).."\"")
    os.exit(1)
end

local tmpfn = os.tmpname()

local fd = assert(io.open(tmpfn,"w"))
fd:write("12345678")
fd:close()

fd,err = io.open(tmpfn,"r+")
if not fd then
    print("NG: io.open tmpfn")
    print(tmpfn,":",err)
    os.exit(1)
end
fd:write("abcd")
fd:close()

for line in io.lines(tmpfn) do
    if line ~= "abcd5678" then
        print("NG: io.write(r+)",line)
        os.exit(1)
    end
    break
end

fd = io.open(tmpfn,"r+")
fd:write("ABCD")
fd:close()

for line in io.lines(tmpfn) do
    if line ~= "ABCD5678" then
        print("NG: io.write(w+)",line)
        os.exit(1)
    end
    break
end
