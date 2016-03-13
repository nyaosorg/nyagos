-- ~\nyagos.ini から、エイリアスや環境変数を読み込むスクリプト

--[[
; ~\nyagos.ini の例(キーの大文字、小文字はあまり区別しないよ)
[alias]
ll=ls -l
lala=ls -al
[env]
path=%path%;c:\hogehoge
]]--
local home=nyagos.env.home or nyagos.env.userprofile
local inipath=nyagos.pathjoin(home,"nyagos.ini")
if nyagos.stat(inipath) then
    local section_key=""
    local all_section={
        alias=nyagos.alias,
        env=nyagos.env
    }
    for line in io.lines(inipath) do
        if string.match(line,"^%s*$") or string.match(line,"^%s*;") then
            goto nextline
        end
        line = string.gsub(line,"%%([^%%]+)%%",function(name)
            return nyagos.env[name]
        end)
        local m=string.match(line,"^%s*%[([^%]]+)%]")
        if m then
            section_key = m
        end
        local key,val = string.match(line,"^%s*([^=]+)=%s*(.*)$")
        if key then
            local sec = all_section[section_key] or {}
            sec[key] = val
            all_section[section_key] = sec 
            -- print(string.format("[%s] %s=%s",section_key,key,val))
        end
        ::nextline::
    end
end
