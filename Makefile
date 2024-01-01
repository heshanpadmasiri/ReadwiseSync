##
# Readwise Sync
#
# @file
# @version 0.1

SRC=$(wildcard *.go)

all: readwisesync

readwisesync: $(SRC)
	go build

install: readwisesync
	mkdir -p ~/bin
	cp readwisesync ~/bin/readwiseSync

unistall:
	rm -f ~/bin/readwiseSync

dist:
	CGO_ENABLED=0 go build -ldflags "-s -w" -o readwiseSync

clean:
	rm -f readwisesync

.PHONEY: dist all install
# end
