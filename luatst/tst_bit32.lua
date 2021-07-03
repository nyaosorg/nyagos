if bit32.band(1,3) ~= 1 then
    os.exit(1)
end

if bit32.bor(1,2) ~= 3 then
    os.exit(1)
end

if bit32.bxor(1,3) ~= 2 then
    os.exit(1)
end
