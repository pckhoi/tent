SHELL=/bin/bash
# A literal space.
space :=
space +=
mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
current_dir := $(dir $(mkfile_path))
dirs_split := $(subst /, , $(current_dir))
go_path := $(subst $(space),/, $(wordlist 1, $(shell expr $(words $(dirs_split)) - 8 ), $(dirs_split)))
tent_src := github.com/pckhoi/tent
postgres_dir := src/$(tent_src)/internal/app/parsers/postgres

bison_to_peg:
	docker run --rm -v $(go_path):/go golang go install $(tent_src)/cmd/bison_to_peg

gen_peg: bison_to_peg
	docker run --rm -v $(go_path):/go golang /bin/bash -c "bison_to_peg \
	$(postgres_dir)/gram.y > \
	$(postgres_dir)/gram.peg"