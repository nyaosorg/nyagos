nyagos.alias.lns = function(args)
    if args and #args == 2 then
        nyagos.shellexecute(
            "runas",
            nyagos.exe,
            string.format([[-c "__ln__ -s '%s' '%s'"]], args[1],args[2])
            )
    else
        print("Usage: ln_s SRCPATH DSTPATH")
    end
end
