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

The easiest way to install `maybe` is to get the [snap package](https://docs.snapcraft.io/core/install):

    sudo snap install maybe


Alternative, compile from source
--------------------------------

To compile the code and make `maybe` system wide accessible:

    go install github.com/thibran/maybe  
    sudo cp $GOPATH/bin/maybe /usr/local/bin


Bash
====

Add to your `~/.bashrc`:

``` bash
function m() {
    if [[ -z $1 ]]; then
        if [[ $PWD != $HOME ]]; then
            cd $HOME
        fi
        return $?
    fi
    d=$(maybe --search $@)
    if [[ $? == 0 ]]; then
        if [[ $d != $PWD ]]; then
            cd $d
        fi
    else
        return 2
    fi
}

function mm() {
    maybe -list $@
}

function cd()
{
    builtin cd $@
    maybe -add $PWD
}
```


Fish Shell
==========

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

(cl-defun maybe (a &optional b &key (fn #'dired))
  (interactive "sMaybe search-query: ")
  (unless (empty-string-p a)
    (let ((dir
           (shell-command-to-string
            (format "maybe -search %s"
                    (mapconcat 'identity (list a b) " ")))))
      (unless (empty-string-p dir)
        (funcall fn dir)))))

(defun maybe-list (a &optional b)
  (unless (empty-string-p a)
    (shell-command-to-string
     (format "maybe -list %s"
             (mapconcat 'identity (list a b) " ")))))

(defun maybe-add-current-folder ()
  "add current folder to the maybe dataset"
  (let ((inhibit-message t)             ; silence echo area output
        (dir (string-remove-prefix "Directory " (pwd))))
    (shell-command
     (format "maybe -add %s" dir))))
```


Eshell
------

``` lisp
(add-hook 'eshell-directory-change-hook #'maybe-add-current-folder)

(defun eshell/m (a &optional b)
  "eshell maybe-search function alias"
  (if (null a)
      (progn (cd "~") ())
    (maybe (symbol-or-string-to-string a)
           (unless (null b) (symbol-or-string-to-string b))
           :fn (lambda (dir) (cd dir) nil))))

(defun eshell/mm (a &optional b)
  "eshell maybe-list alias"
  (unless (null a)
    (maybe-list (symbol-or-string-to-string a)
                (unless (null b) (symbol-or-string-to-string b)))))
```


TODO
====

- write fish completion, using --show with a sub-command
   http://fishshell.com/docs/current/index.html#completion-own
   https://stackoverflow.com/questions/16657803/creating-autocomplete-script-with-sub-commands
   https://github.com/fish-shell/fish-shell/issues/1217#issuecomment-31441757
