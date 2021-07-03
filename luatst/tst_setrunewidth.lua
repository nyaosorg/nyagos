ok,err = nyagos.setrunewidth(100,100)
if not ok then
    print("NG:","return value is not true:",err)
    os.exit(1)
end

ok,err = nyagos.setrunewidth()
if ok then
    print("NG:","did not return parameter error")
    os.exit(1)
end

ok,err = nyagos.setrunewidth(1)
if ok then
    print("NG:","did not return parameter error")
    os.exit(1)
end

ok,err = nyagos.setrunewidth('x')
if ok then
    print("NG:","did not return parameter error")
    os.exit(1)
end

ok,err = nyagos.setrunewidth(100,'x')
if ok then
    print("NG:","did not return parameter error")
    os.exit(1)
end
