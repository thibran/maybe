maybe
=====

version: 0.4

jump to known folder on the command-line

[![asciicast](https://asciinema.org/a/dN7G7dd4GHRiCXMS07CR8GlRg.png)](https://asciinema.org/a/dN7G7dd4GHRiCXMS07CR8GlRg)


Tested on openSUSE Tumbleweed & Ubuntu.

Setup
=====

Fish Shell
----------

Create a fish-function with `funced m` and insert:

```
function m
  if [ "$argv[1]" = "" ]
    # clear
    if [ $PWD != $HOME ]
      cd $HOME
    end
    return $status
  end

  set d (maybe --search $argv)
  if [ $status = 0 ]
    # clear
    if [ $d != $PWD ]
      cd $d
    end
  else
    return 2
  end
end
```

To automatically add visited folders to `maybe`, edit
your `~/.config/fish/config.fish`:

```
function m_on_pwd --on-variable PWD
    maybe -add $PWD
end
```

To list the query results without jumping to the top match:

```
function mm
        maybe -list $argv
end
```

Save both newly created fish functions with:

```
funcsave m
funcsave mm
```

TODO
====

- write fish completion, using --show with a sub-command
   http://fishshell.com/docs/current/index.html#completion-own
   https://stackoverflow.com/questions/16657803/creating-autocomplete-script-with-sub-commands
   https://github.com/fish-shell/fish-shell/issues/1217#issuecomment-31441757