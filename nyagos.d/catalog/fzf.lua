if not nyagos then
    print("This is a script for nyagos not lua.exe")
    os.exit()
end

local fzf = {}
fzf.cmd          =  "fzf.exe"
fzf.args         =  {}
fzf.args.dir     =  ""
fzf.args.cmdhist =  ""
fzf.args.cdhist  =  ""
fzf.args.gitlog  =  "--preview='git show {1}'"

share.fuzzyfinder = fzf

use "fuzzyfinder"
