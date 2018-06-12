nyagos.alias.i320 = function(args)
    assert( nyagos.rawexec( { [0]="cmd",[1]="/c",[2]="dir" } ) )
end

nyagos.alias.i320b = function(args)
    local str = assert( nyagos.raweval( { [0]="cmd",[1]="/c",[2]="dir" } ) )
    print("eval=",str)
end
