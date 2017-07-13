maybe
=====

fish shell function
-------------------

To search, create a fish-function with `funced m` and insert:

```
function m
  if [ "$argv[1]" = "" ]
  clear; and cd $HOME
  return $status
end
  set d (maybe --search $argv)
  if [ "$status" = 0 ]
    clear; and cd $d
  else
    return 1
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
