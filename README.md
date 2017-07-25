maybe
=====

version: 0.3.1

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
