GOGPT
===

**ChatGPT. In your terminal!**

ChatGPT is awesome, and the web interface is great, but as I'm a CLI freak, I wanted ChatGPT be under my finger tips.

The goal is to make the CLI as close to web experience as possible. For that matter I've implemented a REPL. But in case I need quickly ask something, I wanted this to work:

```bash
$ gogpt List sources to learn golang interfaces
```

### Chip in!

Another goal was to learn Go. Being a larva of a gopher I very much welcome contributions. Feel free to open up a PR and/or school me how to go ;)

There are a few things can be done now to improve the REPL, the experimental markdown renderer, in particular.

### REPL with chat capabilities

![](./.github/chat.gif)

### CLI

![](./.github/gogpt.gif)

## Install

For now `gogpt` is NOT distributed via popular package managers such as brew, choco, yum, etc.

Distribution via the main package managers is on the roadmap, but for now you can use one of the following methods (the drawback is updates are not automated, so if you like `gogpt` check back in once in a while):

### curl

```bash
curl https://i.jpillora.com/Nemoden/gogpt! | bash
```

it will install the binary to `/usr/local/bin`

### Releases page

Visit the [releases](https://github.com/nemoden/gogpt/releases) page and download the appropriate archive, unarchive it and drop the binary anywhere you like (ideally on of the directories that are on your `$PATH`.

### Go

```
go install github.com/nemoden/gogpt@latest
```

it will install the binary in your `$GOBIN` directory

## Uninstall

Simply delete the binary and (if they exist)

- config directory located at `~/.config/gopgt`
- the token file `~/.gogpt`

## Roadmap

- [ ] Versioning and release toolings
- [ ] Docker images. Must be super slim!
- [ ] Distribute via brew, etc.
- [ ] Nicer REPL (formatting, colors etc)
- [ ] Enhancements to `gogpt config` command: more configurable options (max tokens, temperature etc)
- [ ] Configurable prompts (i.e. "Write me `<essay>` in the style of `<John Doe>`" where user only needs to type `<essay>` and `<John Doe>`)
- [ ] REPL History (up/down arrows)
- [ ] REPL Ctrl+Enter to enter a new line
- [ ] REPL Option+B / Option+F to jump words, Ctrl+W to delete backward, etc, all the shell goodies to navigate (or vim-like nav?)
- [ ] REPL Completions and commands, i.e. something like what postgresql has, or redis, i.e. being able to type `\config renderer=markdown`, or `\prompt ...`
- [ ] Put all of the above into github issues for tracking and making it easier to collab
