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

if io.type("") == nil then
    print "OK: iotype()==nil"
else
    print("NG: iotype()==\""..io.type("").."\"")
end

if line == "ONELINE" then
    print"OK: read('*l')"
else
    print("NG: read('*l'):[",line,"]")
end
if num == 4 then
    print "OK: read('*n')"
else
    print("NG: read('*n')",num)
end
if crlf == "\n" then
    print "OK: read(1)"
else
    print ("NG: read(1):",crlf)
end
if rest == "AHAHA\nIHIHI\nUFUFU" then
    print "OK: read('*a')"
else
    print("NG: read('*a'):[",rest,"]")
end

local pos,err = fd:seek('set', 0)
if pos == 0 then
  local line = fd:read('*l')
  if line == "ONELINE" then
    print "OK: seek('set',0)"
  else
    print("NG: seek('set',0)", line)
  end
else
  print("NG: seek('set',0)",pos,err)
end

local pos,err = fd:seek('set', 9)
if pos == 9 then
  local line = fd:read('*l')
  if line == '4' then
    print "OK: seek('set',9)==9 read('*l')==4"
  else
    print("NG: seek('set',9)==9 read('*l')==4:",line)
  end
else
  print("NG: seek('set',9)==9:",pos)
end

local pos,err = fd:seek('cur')
if pos == 12  then
  local line = fd:read('*l')
  if line == 'AHAHA' then
    print "OK: seek('cur')==12 read('*l')=='AHAHA'"
  else
    print("NG: seek('cur')==12 read('*l')=='AHAHA'",line)
  end
else
  print("NG: seek('cur')==12:", pos)
end

local pos,err = fd:seek('cur',2)
if pos == 21 then
  local line = fd:read('*l')
  if line == 'IHI' then
    print "OK: seek('cur',2)==21 read('*l')=='IHI'"
  else
    print("NG: seek('cur',2)==21 read('*l')=='IHI'",line)
  end
else
  print("NG: seek('cur',2)==21:",pos)
end

local pos,err = fd:seek('end')
if pos == 31 then
  local empty = fd:read('*a')
  if empty == '' then
    print "OK: seek('end')==31 read('*a')==''"
  else
    print("NG: seek('end')==31 read('*a')==''",empty)
  end
else
  print("NG: seek('end')==31",pos)
end

local size,err = fd:seek('end')
fd:seek('set',0)
rest,err = fd:read(size + 1)
if rest == "ONELINE\n4\nAHAHA\nIHIHI\nUFUFU" then
  print "OK: read(number) number > file size"
else
  print("NG: read(number) number > file size:", rest)
end

fd:seek('set',0)
local all = fd:read('*a')
local before = all:match('.*\n'):gsub('\n','\r\n')
fd:seek('set',#before)
local line = fd:read('*l')
if line == 'UFUFU' then
  print "OK: read('*l') on line without line break"
else
  print("NG: read('*l') on line without line break:",line)
end

fd:seek('end')
local eof = fd:read('*a')
if eof == '' then
  print "OK: read('*a') on EOF"
else
  print("NG: read('*a') on EOF:",eof)
end

local eof = fd:read('*l')
if eof == nil then
  print "OK: read('*l') on EOF"
else
  print("NG: read('*l') on EOF:["..eof..']')
end

local eof = fd:read('*n')
if eof == nil then
  print "OK: read('*n') on EOF"
else
  print("NG: read('*n') on EOF:["..eof..']')
end

local eof = fd:read(1)
if eof == nil then
  print "OK: read(number) on EOF"
else
  print("NG: read('number') on EOF:["..eof.."]")
end

fd:seek('set',0)
local empty = fd:read(0)
if empty == '' then
  print "OK: read(0) on not EOF"
else
  print("NG: read(0) on not EOF:",eof)
end
local _, eof = fd:read('*a', 0)
if eof == nil then
  print "OK: read(0) on EOF"
else
  print("NG: read(0) on EOF:",eof)
end

fd:close()
if io.type(fd) == "closed file" then
    print "OK: iotype()==\"closed file\""
else
    print("NG: iotype()==\""..io.type(fd).."\"")
end

io.stdout:setvbuf("no")
