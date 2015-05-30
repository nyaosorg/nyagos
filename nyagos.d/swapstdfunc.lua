io.getenv = nyagos.getenv
io.setenv = nyagos.setenv

function nyagos.echo(...)
    nyagos.write(...)
    nyagos.write("\n")
end

original_print = print
print = nyagos.echo

function x(s)
    for line in string.gmatch(s,'[^\r\n]+') do
        nyagos.exec(line)
    end
end
