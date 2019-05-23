if not nyagos then
    print("This is a script for nyagos not lua.exe")
    os.exit()
end

nyagos.key.C_o = function(this)
    local word,pos = this:lastword()
    word = string.gsub(word,'"','')
    local wildcard = word.."*"
    local list = nyagos.glob(wildcard)
    if #list == 1 and list[1] == wildcard then
        return
    end
    local dict = {}
    local array = {}
    for _,path in ipairs(list) do
        local index=string.find(path,"[^\\]+$")
        local fname
        if index then
            fname=string.sub(path,index)
        else
            fname=path
        end
        array[1+#array] = fname
        dict[fname] = path
    end
    nyagos.write("\n")
    local result=nyagos.box(array)
    if result then
        result = dict[result]
        if not result then
            result = word
        end
    else
        result = word
    end
    this:call("REPAINT_ON_NEWLINE")
    if string.find(result," ",1,true) then
        if string.find(result,"^~[\\/]") then
            result = '~"'..string.sub(result,2)..'"'
        else
            result = '"'..result..'"'
        end
    end
    assert( this:replacefrom(pos,result) )
end

share.__dump_history = function()
    local uniq={}
    local result={}
    for i=nyagos.gethistory()-1,1,-1 do
        local line = nyagos.gethistory(i)
        if line ~= "" and not uniq[line] then
            result[ #result+1 ] = line
            uniq[line] = true
        end
    end
    return result
end

nyagos.key.C_x = function(this)
    nyagos.write("\nC-x: [r]:command-history, [h]:cd-history, [g]:git-revision\n")
    local ch = nyagos.getkey()
    local c = string.lower(string.char(ch))
    local result
    if c == 'r' or ch == bit32.band(string.byte('r'),0x1F) then
        result = nyagos.box(share.__dump_history())
    elseif c == 'h' or ch == bit32.band(string.byte('h') , 0x1F) then
        result = nyagos.eval('cd --history | box')
        if string.find(result,' ') then
            result = '"'..result..'"'
        end
    elseif c == 'g' or ch == bit32.band(string.byte('g'),0x1F) then
        result = nyagos.eval('git log --pretty="format:%h %s" | box')
        result = string.match(result,"^%S+") or ""
    end
    this:call("REPAINT_ON_NEWLINE")
    return result
end

nyagos.key.M_r = function(this)
    nyagos.write("\n")
    local result = nyagos.box(share.__dump_history())
    this:call("REPAINT_ON_NEWLINE")
    if string.find(result,' ') then
        result = '"'..result..'"'
    end
    return result
end

nyagos.key.M_h = function(this)
    nyagos.write("\n")
    local result = nyagos.eval('cd --history | box')
    this:call("REPAINT_ON_NEWLINE")
    if string.find(result,' ') then
        result = '"'..result..'"'
    end
    return result
end

nyagos.key.M_g = function(this)
    nyagos.write("\n")
    local result = nyagos.eval('git log --pretty="format:%h %s" | box')
    this:call("REPAINT_ON_NEWLINE")
    return string.match(result,"^%S+") or ""
end

nyagos.key["M-o"] = function(this)
    local spacecut = false
    if string.match(this.text," $") then
        this:call("BACKWARD_DELETE_CHAR")
        spacecut = true
    end
    local path,pos = this:lastword()
    if not string.match(path,"%.[Ll][Nn][Kk]$") then
        return
    end
    path = string.gsub(path,'"','')
    path = string.gsub(path,"/","\\")
    path = string.gsub(path,"^~",os.getenv("USERPROFILE"))

    local wsh,err = nyagos.create_object("WScript.Shell")
    if wsh then
        local shortcut = wsh:CreateShortCut(path)
        if shortcut then
            local newpath = shortcut:_get("TargetPath")
            if newpath then
                local isDir = false
                local fso = nyagos.create_object("Scripting.FileSystemObject")
                if fso then
                    if fso:FolderExists(newpath) then
                        isDir = true
                    end
                    fso:_release()
                end
                if string.find(newpath," ") then
                    if isDir then
                        newpath = '"'..newpath..'\\'
                    else
                        newpath = '"'..newpath..'"'
                        if spacecut then
                            newpath = newpath .. ' '
                        end
                    end
                elseif isDir then
                    newpath = newpath .. '\\'
                elseif spacecut then
                    newpath = newpath .. ' '
                end
                if string.len(newpath) > 0 then
                    this:replacefrom(pos,newpath)
                end
            end
            shortcut:_release()
        end
        wsh:_release()
    end
end
