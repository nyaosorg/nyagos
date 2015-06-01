nyagos.on_command_not_found = function(args)
    nyagos.writerr(args[0]..": コマンドではない。\n")
    return true
end
-- vim:set fenc=utf8 --
