local sample = "あいうえお"
local r = utf8.len(sample)
print( "len('"..sample.."')==",r)
assert( r == 5)
local r1,r2 = utf8.len(sample,3,-4)
print(r1,r2)
assert( r1 == nil and r2 == 3)

local r1,r2 = utf8.len(sample,1,-1)
print(r1,r2)
assert( r1 == 5)

local r1,r2 = utf8.len(sample,1,-2)
print(r1,r2)
assert( r1 == 5)

assert( utf8.len(sample,1,-3) == 5)
assert( utf8.len(sample,1,-4) == 4)
assert( utf8.len(sample,1,-5) == 4)

local r1,r2 = utf8.len(sample,2)
assert( r1 == nil and r2 == 2 )
local r1,r2 = utf8.len(sample,3)
assert( r1 == nil and r2 == 3 )
