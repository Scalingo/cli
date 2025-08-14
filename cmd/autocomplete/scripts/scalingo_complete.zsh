autoload -U compinit && compinit
autoload -U bashcompinit && bashcompinit

#!/usr/bin/env zsh

_cli_bash_autocomplete() {
     local cur prev opts base
     COMPREPLY=()
     cur="${COMP_WORDS[COMP_CWORD]}"
     prev="${COMP_WORDS[COMP_CWORD-1]}"
     opts=$( ${COMP_WORDS[@]:0:$COMP_CWORD} --generate-shell-completion )
     COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
     return 0
 }

complete -F _cli_bash_autocomplete scalingo
