# jjmp

Bookmarks for the CLI.

<img width="658" alt="j" src="https://github.com/user-attachments/assets/1762b829-1e62-4744-8fbb-e3f2ff241947">

## Philosophy

Imagine having a short list of bookmarks available in your CLI. You have one for your current project, one for your other current project, one for your Neovim setup, etc. Imagine being able to jump between them with one simple command.

Now imagine you only have 10 bookmark slots, so you have to be judicious as to what you bookmark, and how you maintain the list, so that it continues being useful.

That's the idea behind `jjmp`.

### Why not use symlinks?

Symlinks tend to be permanent. They are easy to set up, but equally easy to never delete, so you end up with a lot of clutter in your home directory. A limited list of 10 bookmarks forces you to actively prune it, and make decisions about to what should be on it, based on how often you use it.

### Why not use a shell script to do all this?

I have done this [here](https://gist.github.com/maciakl/b7f65bf40a1a78c06e6b0b058d76234f).

This implementation had a couple of flaws:

1. I needed separate script for Powershell and Zsh with slightly different syntax and logic
1. The scripts were dependent on Skate and Gum executables for key-value store and UI respectively
1. The Skate and Gum updates kept breaking the them

By virtue of being a go program `jjmp` is platform agnostic and does not depend on any external executables. It also uses a simpler key-value store that is unique to it and less likely to get clobbered by updates.

### Why not use zoxide instead?

[zoxide](https://github.com/ajeetdsouza/zoxide) is a great tool, but it has a slightly different focus. It is designed to entirely replace your `cd` command in your workflow. I just wanted a manually curated bookmark list. 

## Installing

 Install via go:
 
    go install github.com/maciakl/jjmp@latest
 
 On Windows, this tool is also distributed via `scoop` (see [scoop.sh](https://scoop.sh)).

 First, you need to add my bucket:

    scoop bucket add maciak https://github.com/maciakl/bucket
    scoop update

 Next simply run:
 
    scoop install jjmp

If you don't want to use `scoop` you can simply download the executable from the release page and extract it somewhere in your path.

## Setup

Because an executable cannot change the working directory of a parent shell you need to wrap the `jjmp` command in a shell function.

Poweshell:
```powershell
function j { jjmp.exe $args | cd }
```

Bash & Zsh:
```zsh
function j { cd "$(jjmp $@)" }
```

## Usage

Onc you wrapped the `jjmp` command in a shell function `j` you can use it like this:

- `j` - will list all bookmarks
- `j <0-9>` - will change the working directory to the bookmarked path
- `j set <0-9>` - will bookmark the current working directory
- `j delete <0-9>` - will delete the chosen bookmark

Bookmarks are stored in `~/.jjmpdb` file.
