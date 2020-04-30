#!/bin/sh

zcat Core_v5.2_modified.txt.gz >/tmp/core.txt

perl parse.pl /tmp/core.txt 0 >parsed.txt
perl parse.pl /tmp/core.txt 1 >parsed_events.txt

mkdir -p output/
rm output/*
perl generate.pl parsed.txt 1 >output/link_control.go
perl generate.pl parsed.txt 2 >output/link_policy.go
perl generate.pl parsed.txt 3 >output/baseband.go
perl generate.pl parsed.txt 4 >output/informational.go
perl generate.pl parsed.txt 5 >output/status.go
perl generate.pl parsed.txt 6 >output/testing.go
perl generate.pl parsed_events.txt 7 >output/events.go
perl generate.pl parsed.txt 8 >output/le.go

rm parsed.txt
rm parsed_events.txt
