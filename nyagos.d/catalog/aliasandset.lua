if not nyagos then
    print("This is a script for nyagos not lua.exe")
    os.exit()
end
local function expand(text)
    local result,_ = string.gsub(text,"%%([%w_%(%)]+)%%",function(w)
        return nyagos.getenv(w)
    end)
    return result
end

local function set_(f,equation,expand)
    if type(equation) == 'table' then
        for left,right in pairs(equation) do
            f(left,expand(right))
        end
        return true
    end
    local pluspos=string.find(equation,"+=",1,true)
    if pluspos and pluspos > 0 then
        local left=string.sub(equation,1,pluspos-1)
        equation = string.format("%s=%s;%%%s%%",
                        left,string.sub(equation,pluspos+2),left)
    end
    local pos=string.find(equation,"=",1,true)
    if pos then
        local left=string.sub(equation,1,pos-1)
        local right=string.sub(equation,pos+1)
        f( left , expand(right) )
        return true
    end
    return false,(equation .. ': invalid format')
end

function set(equation) 
    set_(nyagos.setenv,equation,expand)
end
function alias(equation)
    set_(nyagos.setalias,equation,function(x) return x end)
end
