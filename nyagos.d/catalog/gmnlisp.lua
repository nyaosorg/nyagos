if not nyagos then
    print("This is a script for nyagos not lua.exe")
    os.exit()
end

local orgfilter = nyagos.filter
nyagos.filter = function(cmdline)
    if string.sub(cmdline,1,1) == '(' then
        cmdline = string.gsub(cmdline,'"','\\"')
        cmdline = string.format('gmnlisp -e "(format t \\\"~a\\\" %s)"', cmdline)
        nyagos.exec(cmdline)
        print()
        return ""
    end
    cmdline = string.gsub(cmdline,'@%b()', function(code)
        code = string.sub(code,2)
        code = string.gsub(code,'"','\\"')
        code = string.format('gmnlisp -e "(format t \\\"~a\\\" %s)"', code)
        return nyagos.eval(code)
    end)
    if orgfilter then
        local cmdline_ = orgfilter(cmdline)
        if cmdline_ then
            cmdline = cmdline_
        end
    end
    return cmdline
end
