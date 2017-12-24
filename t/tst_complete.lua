-- How to test:
--    `lua_f tst_complete.lua`
--    and see result of completion.

nyagos.completion_hook = function(c)
    if c.field[1] == "svn" then
        c.list[ #c.list+1 ] = "commit"
        c.list[ #c.list+1 ] = "update"
        c.list[ #c.list+1 ] = "ls"
    end
    nyagos.msgbox("leftStr="..c.left)
    return c.list
end
