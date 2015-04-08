# Go Challenge 1

This project is my solution for the first [Go Challenge](http://golang-challenge.com/go-challenge1/).

Its goal is to read in a binary file and decode its contents to output
a drum track pattern.

I've been playing with Go for a while and really value getting back to a language
with a type system that still bears the feel of a lightweight and fast development
workflow.

# Running

There isn't too much to know. Run `go get github.com/thedahv/gochallenge1` to download
the project into your computer.

Run `go test ./...` in the project folder to run all tests.

# Lessons Learned

This was also a first dive for me into low-level programming and parsing since college.
Here are a few things I picked up as I worked on this challenge.

## Go Standard Library is Great

The Go standard library is fantastic and really offers many tools for solving problems
like these out of the box.

Given more time, I would have spent more time reading and understanding the
[`binary.Read`](http://golang.org/pkg/encoding/binary/#Read) function and related
helpers to simplify a lot of my byte counting and error checking.

## Go Testing is Fun!

There are varying opinions on how elegant or verbose Go testing might be.
I think Go testing is a great asset for the language, and tools like
[GoConvey](http://goconvey.co/) are great for adding an elegant test runner
and analysis tool to your project.

## The Go Community is _Really_ Helpful

This was one of the first times I interacted with a software community on Slack
as opposed to IRC. However, the people there are encouraging, smart, and helpful.

I didn't find out about the Slack community until the evening before the due date,
but I've been involved there ever since.
