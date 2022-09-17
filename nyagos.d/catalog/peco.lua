if not nyagos then
    print("This is a script for nyagos not lua.exe")
    os.exit()
end

local peco = {}
peco.cmd          =  "peco.exe"
peco.args         =  {}
peco.args.dir     =  ""
peco.args.cmdhist =  ""
peco.args.cdhist  =  ""
peco.args.gitlog  =  ""

share.fuzzyfinder = peco

use "fuzzyfinder"
