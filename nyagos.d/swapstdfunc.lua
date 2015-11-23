io.getenv = nyagos.getenv
io.setenv = nyagos.setenv

print = function(...)
    nyagos.write(...)
    nyagos.write("\n")
end

function x(s)
    for line in string.gmatch(s,'[^\r\n]+') do
        nyagos.exec(line)
    end
end
