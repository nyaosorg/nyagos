share.foo = {}
share.foo.ahaha = "ahaha"
local x = share.foo
print(share.foo.ahaha or "(nil)")
if share.foo.ahaha == "ahaha" then
    print("-> OK")
else
    print("-> NG")
end
share.foo = {}
x.ihihi = "fooo"
print(share.foo.ihihi or "(nil)")
if share.foo.ihihi == nil then
    print("-> OK")
else
    print("-> NG")
end
