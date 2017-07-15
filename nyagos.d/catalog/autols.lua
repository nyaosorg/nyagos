if not nyagos then
    print("This is a script for nyagos not lua.exe")
    os.exit()
end

share._autols = function(func, cmd, args)
  local status, err
  if func then
    status, err = func(args)
  else
    args[0] = cmd
    status, err = nyagos.exec(args)
  end
  if not status then
    nyagos.exec('ls')
  end
  return status, err
end

share.org_autols_cd = nyagos.alias.cd
nyagos.alias.cd = function(args)
  return share._autols(share.org_autols_cd, '__cd__', args)
end

share.org_autols_pushd = nyagos.alias.pushd
nyagos.alias.pushd = function(args)
  return share._autols(share.org_autols_pushd, '__pushd__', args)
end

share.org_autols_popd = nyagos.alias.popd
nyagos.alias.popd = function(args)
  return share._autols(share.org_autols_popd, '__popd__', args)
end
