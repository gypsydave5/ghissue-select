# ghissue-select

These docs are WIP :)

This is based on Tamara Jordan's [coauthor-select](https://github.com/tamj0rd2/coauthor-select) tool.

And when I say "based on", basically ripped off and mauled. Sorry Tam.

- **cmd/select** - allows you to select coauthors from a list. Add this to your prepare-commit-msg hook

## How to use these tools as git hooks (if you're using go modules):

1. Add [tools.go](./example/tools.go) to your project
2. Add `require github.com/gypsydave5/ghissue-select v0.1.0` to your `go.mod`
3. Run `go mod tidy && go mod vendor`
4. Create a hooks folder `mkdir .hooks` in your project and enable it as the git hooks folder `git config core.hooksPath .hooks`
5. Copy /examples/.hooks to your repo and make all files executable

## Specifying pairs via the command line

1. Commit as you usually do
2. You'll be prompted to enter a GitHub Issue ID or be given the option to choose the last ID you used. This is enabled by the prepare-commit-msg hook.
3. You'll be warned if you're trying to commit to the trunk without specifying a pair

## Configuration

### cmd/select

Check [here](./cmd/select/main) for defaults and the latest documentation

- `--issueFile` - the path to your .ghissue file
- `--interactive` - set this to false if you're using a non-interactive console

