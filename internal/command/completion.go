package command

import "fmt"

// runCompletion handles the completion subcommand.
func runCompletion(args []string, deps *Deps) error {
	if len(args) == 0 {
		printCompletionUsage(deps)
		return nil
	}

	switch args[0] {
	case "bash":
		fmt.Fprint(deps.Stdout, bashCompletion)
		return nil
	case "zsh":
		fmt.Fprint(deps.Stdout, zshCompletion)
		return nil
	case "-help", "--help", "-h":
		printCompletionUsage(deps)
		return nil
	default:
		return fmt.Errorf("completion: unsupported shell %q (use bash or zsh)", args[0])
	}
}

// printCompletionUsage writes completion command help text.
func printCompletionUsage(deps *Deps) {
	fmt.Fprint(deps.Stdout, `Usage: lcli completion <shell>

Shells:
  bash    Generate bash completion script
  zsh     Generate zsh completion script

Install:
  lcli completion bash > /etc/bash_completion.d/lcli
  lcli completion zsh > "${fpath[1]}/_lcli"
`)
}

const bashCompletion = `# bash completion for lcli -*- shell-script -*-

_lcli() {
    local cur prev commands
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    commands="auth config profile post comment reaction media org analytics completion version help"

    case "${prev}" in
        lcli)
            COMPREPLY=( $(compgen -W "${commands}" -- "${cur}") )
            return 0
            ;;
        auth)
            COMPREPLY=( $(compgen -W "login logout status" -- "${cur}") )
            return 0
            ;;
        config)
            COMPREPLY=( $(compgen -W "setup" -- "${cur}") )
            return 0
            ;;
        profile)
            COMPREPLY=( $(compgen -W "me view" -- "${cur}") )
            return 0
            ;;
        post)
            COMPREPLY=( $(compgen -W "create list get delete" -- "${cur}") )
            return 0
            ;;
        comment)
            COMPREPLY=( $(compgen -W "create list delete" -- "${cur}") )
            return 0
            ;;
        reaction)
            COMPREPLY=( $(compgen -W "like unlike list" -- "${cur}") )
            return 0
            ;;
        media)
            COMPREPLY=( $(compgen -W "upload" -- "${cur}") )
            return 0
            ;;
        org)
            COMPREPLY=( $(compgen -W "info followers stats" -- "${cur}") )
            return 0
            ;;
        analytics)
            COMPREPLY=( $(compgen -W "post views" -- "${cur}") )
            return 0
            ;;
        completion)
            COMPREPLY=( $(compgen -W "bash zsh" -- "${cur}") )
            return 0
            ;;
    esac
}

complete -F _lcli lcli
`

const zshCompletion = `#compdef lcli

_lcli() {
    local -a commands
    commands=(
        'auth:Authenticate with LinkedIn'
        'config:Configure client credentials'
        'profile:View LinkedIn profiles'
        'post:Create, list, and manage posts'
        'comment:Manage comments on posts'
        'reaction:Like and react to posts'
        'media:Upload images and videos'
        'org:Manage organization pages'
        'analytics:View post and profile analytics'
        'completion:Generate shell completions'
        'version:Print version information'
        'help:Show usage information'
    )

    _arguments -C \
        '1:command:->command' \
        '*::arg:->args'

    case $state in
        command)
            _describe -t commands 'lcli command' commands
            ;;
        args)
            case $words[1] in
                auth)
                    _values 'subcommand' 'login[Authenticate via OAuth]' 'logout[Remove stored credentials]' 'status[Show auth status]'
                    ;;
                config)
                    _values 'subcommand' 'setup[Configure credentials]'
                    ;;
                profile)
                    _values 'subcommand' 'me[Show your profile]' 'view[View another profile]'
                    ;;
                post)
                    _values 'subcommand' 'create[Create a new post]' 'list[List recent posts]' 'get[Get a single post]' 'delete[Delete a post]'
                    ;;
                comment)
                    _values 'subcommand' 'create[Add a comment]' 'list[List comments]' 'delete[Delete a comment]'
                    ;;
                reaction)
                    _values 'subcommand' 'like[React to a post]' 'unlike[Remove a reaction]' 'list[List reactions]'
                    ;;
                media)
                    _values 'subcommand' 'upload[Upload an image or video]'
                    ;;
                org)
                    _values 'subcommand' 'info[Get organization info]' 'followers[Get follower stats]' 'stats[Get page stats]'
                    ;;
                analytics)
                    _values 'subcommand' 'post[View post analytics]' 'views[View profile views]'
                    ;;
                completion)
                    _values 'shell' 'bash[Generate bash completions]' 'zsh[Generate zsh completions]'
                    ;;
            esac
            ;;
    esac
}

_lcli
`
