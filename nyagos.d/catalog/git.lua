share.git = {}

-- setup branch detector
local branchdetect = function()
  local gitbranches = {}
  local gitbranch_tmp = nyagos.eval('git for-each-ref  --format="%(refname:short)" refs/heads/ 2> nul')
  for line in gitbranch_tmp:gmatch('[^\n]+') do
    table.insert(gitbranches,line)
  end
  return gitbranches
end

-- subcommands
local gitsubcommands={}
-- keyword
gitsubcommands["bisect"]={"start", "bad", "good", "skip", "reset", "visualize", "replay", "log", "run"}
gitsubcommands["notes"]={"add", "append", "copy", "edit", "list", "prune", "remove", "show"}
gitsubcommands["reflog"]={"show", "delete", "expire"}
gitsubcommands["rerere"]={"clear", "forget", "diff", "remaining", "status", "gc"}
gitsubcommands["stash"]={"save", "list", "show", "apply", "clear", "drop", "pop", "create", "branch"}
gitsubcommands["submodule"]={"add", "status", "init", "deinit", "update", "summary", "foreach", "sync"}
gitsubcommands["svn"]={"init", "fetch", "clone", "rebase", "dcommit", "log", "find-rev", "set-tree", "commit-diff", "info", "create-ignore", "propget", "proplist", "show-ignore", "show-externals", "branch", "tag", "blame", "migrate", "mkdirs", "reset", "gc"}
gitsubcommands["worktree"]={"add", "list", "lock", "prune", "unlock"}

-- branch
gitsubcommands["checkout"]=branchdetect
gitsubcommands["reset"]=branchdetect
gitsubcommands["merge"]=branchdetect
gitsubcommands["rebase"]=branchdetect

share.git.subcommand=gitsubcommands
share.git.branch = branchdetect

if share.maincmds then
  if share.maincmds["git"] then
    -- git command complementation exists.
    local maincmds = share.maincmds

    -- build
    for key, cmds in pairs(gitsubcommands) do
      local gitcommand="git "..key
      maincmds[gitcommand]=cmds
    end

    -- replace
    share.maincmds = maincmds
  end
end
