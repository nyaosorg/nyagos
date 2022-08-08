share.foo = {}
share.foo.ahaha = "ahaha"
local x = share.foo
-- print(share.foo.ahaha or "(nil)")
assert(share.foo.ahaha == "ahaha")
share.foo = {}
x.ihihi = "fooo"
-- print(share.foo.ihihi or "(nil)")
assert(share.foo.ihihi == nil)
