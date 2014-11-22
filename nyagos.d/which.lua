nyagos.alias("which",function(args)
    local result=nyagos.which(args[1])
    if result then
        nyagos.write( result..'\n')
    end
end)
