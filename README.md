maybe
=====

version: 0.5.0

jump to known folder on the command-line

[![asciicast](https://asciinema.org/a/dN7G7dd4GHRiCXMS07CR8GlRg.png)](https://asciinema.org/a/dN7G7dd4GHRiCXMS07CR8GlRg)


Tested on openSUSE Tumbleweed, Ubuntu and FreeBSD.


Flags
-----

    -init
          scan $HOME and add folders (six folder-level deep)
    -list string
          list results for keyword
    -datadir string
          (default $HOME/.local/share/maybe)
    -max-entries int
          maximum unique path-entries (default 10000)
    -add string
          add path to index
    -search string
          search for keyword
    -v    verbose
    -version
          print maybe version


Install
=======

Snap Package
------------

The easiest way to install Maybe is to get the [snap](https://docs.snapcraft.io/core/install) package:

    sudo snap install maybe


Alternative, compile from source
--------------------------------

1. Compile the code with:

    git clone https://github.com/thibran/maybe.git
    cd maybe
    go build

2. Make the binary system wide accessible:

    sudo cp maybe /usr/local/bin


Fish Shell
----------

Create `/etc/fish/functions/m.fish` and insert:

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


To list results without jumping to the top match create `/etc/fish/functions/mm.fish` and insert:

```
function mm
        maybe -list $argv
end
```


To automatically add visited folders to `maybe`, create/edit `/etc/fish/config.fish`:

```
function m_on_pwd --on-variable PWD
    maybe -add $PWD
end
```


Emacs
=====

``` lisp
(defun empty-string-p (str)
  (or (null str) (string= "" str)))

(cl-defun maybe (query &optional (fn #'dired) )
  (interactive "sMaybe search-query: ")
  (let ((dir
          (shell-command-to-string
            (format "maybe -search %s" query))))
    (unless (empty-string-p query)
      (funcall fn dir))))

(defun maybe-list (query)
  "list results for query"
  (unless (empty-string-p query)
    (shell-command-to-string (format "maybe -list %s" query))))

(defun maybe-add-current-folder ()
  "add current folder to the maybe dataset"
  (let ((inhibit-message t) ; silence echo area output
	(dir (string-remove-prefix "Directory " (pwd))))
    (shell-command
     (format "maybe -add %s" dir))))
```


Eshell
------

``` lisp
(add-hook 'eshell-directory-change-hook #'maybe-add-current-folder)

(defun eshell/m (&rest q)
  "eshell maybe-search function alias"
  (if (null q)
      (progn (cd "~") ())
    (maybe (mapconcat #'symbol-or-string-to-string q " ")
           (lambda (dir) (cd dir) nil))))

(defun eshell/mm (&rest q)
  "eshell maybe-list alias"
  (unless (null q)
    (maybe-list (mapconcat #'symbol-or-string-to-string q " "))))
```


TODO
====

- write fish completion, using --show with a sub-command
   http://fishshell.com/docs/current/index.html#completion-own
   https://stackoverflow.com/questions/16657803/creating-autocomplete-script-with-sub-commands
   https://github.com/fish-shell/fish-shell/issues/1217#issuecomment-31441757