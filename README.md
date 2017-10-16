# tmux-decerative-launcher
> Write your ideal tmux state in a yaml file and launch it from your shell

## State of this project

Consider this project a small spike whether such a launcher is desirable. I will dogfood it heavily and then decide on whether this is the correct solution or not.

Remember, a spike is dirty. This will be mostly some bash hacking with little to no tests. If I deem this project worthy after an initial test phase I might turn it into an actual project with real quality standards.

### Tech Stack

~I want to give the bash yaml parser [yq](https://github.com/kislyuk/yq) a shot, so this is a bash only project for now.~
Parsing yaml files with `yq` turned out to be ok-ish. However, iterating over arrays was not much fun, so I switched to Golang after the very first commit or so. There are still no tests for now, though.

## Roadmap / Features

* [x] create window with splits based on yml config 
* [x] name window according to template
* [x] run command per pane 
* [ ] specify size for splits
* [ ] don't recreate window when window with name already exists

## Known issues
Right now we're splitting the split. So the more splits you add, they get exponentially smaller. Try with four instead of two.
