--- Do 'lua_f THIS_SCRIPT'

function assert_error(value, error_message)
    assert(value == nil)
    assert(error_message ~= nil)
end


assert(utf8.offset("", 1) == 1)
assert(utf8.offset("", 2) == nil)
assert(utf8.offset("", -1) == nil)
assert(utf8.offset("", 1, 1) == 1)
assert_error(utf8.offset("", 1, 2))
assert_error(utf8.offset("", 1, -1))


local test_string = "aあ𠮷"
-- utf8.offset(s, n, i)
--  s  a|あ   |𠮷     |
--  i  1 2 3 4 5 6 7 8 9
-- -i  8 7 6 5 4 3 2 1
--  n  1 2     3       4
-- -n  3 2     1

-- normal case
assert(utf8.offset(test_string, 1) == 1)
assert(utf8.offset(test_string, 2) == 2)
assert(utf8.offset(test_string, 3) == 5)
assert(utf8.offset(test_string, 4) == 9)
assert(utf8.offset(test_string, 5) == nil)

-- n is negative
assert(utf8.offset(test_string, -4) == nil)
assert(utf8.offset(test_string, -3) == 1)
assert(utf8.offset(test_string, -2) == 2)
assert(utf8.offset(test_string, -1) == 5)

-- i is specified
assert(utf8.offset(test_string, 2, 5) == 9)
assert(utf8.offset(test_string, -1, 5) == 2)
assert(utf8.offset(test_string, 2, 1) == 2)
assert(utf8.offset(test_string, -2, 9) == 2)
assert(utf8.offset(test_string, 2, -4) == 9)
assert(utf8.offset(test_string, -1, -4) == 2)
assert(utf8.offset(test_string, 2, -8) == 2)

-- special case: n is 0
assert(utf8.offset(test_string, 0, 1) == 1)
assert(utf8.offset(test_string, 0, 4) == 2)
assert(utf8.offset(test_string, 0, 6) == 5)
assert(utf8.offset(test_string, 0, 9) == 9)

-- Error case: position out of range
assert_error(utf8.offset(test_string, 1, 0))
assert_error(utf8.offset(test_string, 1, 10))
assert_error(utf8.offset(test_string, 1, -9))

-- Error case: invalid start position
assert_error(utf8.offset(test_string, 1, 3))
assert_error(utf8.offset(test_string, 1, -2))
