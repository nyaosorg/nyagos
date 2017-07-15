if not nyagos then
    print("This is a script for nyagos not lua.exe")
    os.exit()
end

share.org_autocd_argsfilter = nyagos.argsfilter
nyagos.argsfilter = function(args)
  if share.org_autocd_argsfilter then
    local args_ = share.org_autocd_argsfilter(args)
    if args_ then
      args = args_
    end
  end
  if nyagos.which(args[0]) then
    return
  end
  local stat = nyagos.stat(args[0])
  if not stat or not stat.isdir then
    return
  end
  local newargs = {[0] = 'cd'}
  for i = 0, #args do
    newargs[#newargs + 1] = args[i]
  end
  return newargs
end
