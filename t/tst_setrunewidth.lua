ok,err = nyagos.setrunewidth(100,100)
if not ok then
    print("NG:","return value is not true:",err)
else
    print("OK:","normal case")
end

ok,err = nyagos.setrunewidth()
if ok then
    print("NG:","did not return parameter error")
else
    print("OK:", err)
end

ok,err = nyagos.setrunewidth(1)
if ok then
    print("NG:","did not return parameter error")
else
    print("OK: " , err)
end

ok,err = nyagos.setrunewidth('x')
if ok then
    print("NG:","did not return parameter error")
else
    print("OK:" , err)
end

ok,err = nyagos.setrunewidth(100,'x')
if ok then
    print("NG:","did not return parameter error")
else
    print("OK:" , err)
end
