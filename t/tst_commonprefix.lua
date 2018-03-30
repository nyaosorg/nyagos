
val,err = nyagos.commonprefix{ "ABC123" , "ABCCCC" }
if val then
    if val == "ABC" then
        print("OK:",val)
    else
        print("NG:",val)
    end
else
    print("NG:",err)
end

val,err = nyagos.commonprefix()
if not val then
    print("OK:",err)
else
    print("NG:",val)
end

val,err = nyagos.commonprefix(1)
if not val then
    print("OK:",err)
else
    print("NG:",val)
end
