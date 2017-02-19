-- How to test:
--    `lua_f tst_complete.lua`
--    and see result of completion.

nyagos.completion_hook = function(c)
    c.list[ #c.list ] = c.rawword.."(".. c.pos .. ")"
    return c.list
end
