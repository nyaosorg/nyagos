-- For example, on ~/.nyagos
--   local runall = require("runall")
--   runall(nyagos.env.NYAGOS_PRIVATE_DIRECTORIES)
--
--   require("runall")( nyagos.env.NYAGOS_PRIVATE_DIRECTORIES )

return function(dirs)
    for dir in string.gmatch(dirs,"[^;]+") do
        if string.match(dir,"^%~[/\\]") then
            dir = (nyagos.env.HOME or nyagos.env.USERPROFILE or "~") .. string.sub(dir,2)
        end
        local scripts = nyagos.glob(nyagos.pathjoin(dir,"*.lua"))
        if scripts then
            for i = 1, #scripts do
                local chunk,err=loadfile(scripts[i])
                if not chunk then
                    io.stderr:write(scripts[i], ": ", err, "\n")
                else
                    local ok,err = pcall(chunk)
                    if not ok then
                        io.stderr:write(scripts[i], ": ", err, "\n")
                    end
                end
            end
        end
    end
end
