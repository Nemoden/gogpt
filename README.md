GOGPT
===

**ChatGPT. In your terminal!**

ChatGPT is awesome, and the web interface is great, but as I'm a CLI freak, I wanted ChatGPT be under my finger tips.

The goal is to make the CLI as close to web experience as possible. For that matter I've implemented a REPL. But in case I need quickly ask something, I wanted this to work:

```bash
$ gogpt List sources to learn golang interfaces
```

#### Chip in!

Another goal was to learn Go. Being a larva of a gopher I very much welcome contributions. Feel free to open up a PR and/or school me how to go ;)

There are a few things can be done now to improve the REPL. The experimental markdown renderer, in particular. But the most wanted feature is sessions.

#### CLI

![](./.github/gogpt.gif)

#### REPL

![](./.github/repl.gif)

### Roadmap

- [ ] Versioning and release toolings
- [ ] Distribute via brew, etc.
- [ ] Nicer REPL (formatting, colors etc)
- [ ] Enhancements to `gogpt config` command: more configurable options (max tokens, temperature etc)
- [ ] Configurable prompts (i.e. "Write me <foo> in the style of <bar>" where user only needs to type <foo> and <bar>)
- [ ] REPL Sessions (make pass previous context in the next request, i.e. ask something, hit enter, ask to elaborate on previous response)
- [ ] REPL History (up/down arrows)
- [ ] REPL Ctrl+Enter to enter a new line
- [ ] REPL Option+B / Option+F to jump words, Ctrl+W to delete backward, etc, all the shell goodies to navigate (or vim-like nav?)
- [ ] REPL Completions and commands, i.e. something like what postgresql has, or redis, i.e. being able to type `\config renderer=markdown`, or `\prompt ...`
