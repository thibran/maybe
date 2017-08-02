maybe
=====

version: 0.3.1

TODO
====

- jut too long Repo.List filename - ADD flag like: cutLong bool
- -t  by time only
- -r  remove current dir from index
- write fish completion, using --show with a sub-command
   http://fishshell.com/docs/current/index.html#completion-own
   https://stackoverflow.com/questions/16657803/creating-autocomplete-script-with-sub-commands
   https://github.com/fish-shell/fish-shell/issues/1217#issuecomment-31441757
- maybe replace time rating with: fewer seconds from now > better
   if a time value is not present, add penalty

fish shell function
-------------------

To search, create a fish-function with `funced m` and insert:

```
function m
  if [ "$argv[1]" = "" ]
    clear
    if [ $PWD != $HOME ]
      cd $HOME
    end
    return $status
  end

  set d (maybe --search $argv)
  if [ $status = 0 ]
    clear
    if [ $d != $PWD ]
      cd $d
    end
  else
    return 2
  end
end
```

To automatically add visited folders edit
`~/.config/fish/config.fish` and insert:

```
function m_on_pwd --on-variable PWD
    maybe -add $PWD
end
```

To check which other query results are known to maybe:

```
function mm
        maybe -list $argv
end
```
