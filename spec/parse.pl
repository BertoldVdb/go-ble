#!/usr/bin/perl

use strict;
my $active = 0;
my $state = 0;
open (FILE, @ARGV[0]) or die "File not found";

my $event = int(@ARGV[1]);

while(<FILE>){
    chomp $_;
    chomp $_;

    if ($_ eq "For the Link Control commands, the OGF is defined as 0x01."){
        $active=1;
    }
    last if ($_ eq "Commands, events, and configuration parameters in this section were in prior");
    if ($active) {
        if (($event == 0 && $state == 0 && $_ =~ m/^7\.(.*?)\.(.*?) (.*) (Command|command).?$/) || 
            ($event == 1 && $_ =~ m/^7\.(.*?)\.(.*?) (.*) (Event|event)$/)){
            if ($event == 0 || $1 eq "7" && $3 ne "LE Meta"){
                if ($event == 1 && $state == 4){
                    print " |5|TERM\n";
                }
                print "c|$1|$2|$3|";
                $state = 1;
            }
        }
        if ($state == 1 && $_ =~ m/^\sHCI_(.*?)\s0x(.*?)\s/){
            print "0x$2\n";
            $state = 2;
        }
        if ($state == 2 &&
                (($event == 0 && $_ =~ m/^Command parameters:/) || 
                ($event == 1 && $_ =~ m/^Event (parameters|Parameters):/))){
            $state = 3+$event;
        }
        if ($state == 3 || $state == 4){
            if ($_ =~ m/^(.*?):\s\s+(.*?)\soctet/){
                my $size = $2;
                if (substr($size,0,5) eq "Size:"){
                    $size=substr($size, 6, 99999);
                }
                my $name = $1;

                if($size =~ m/(.*)\s(.*)\s(.*?)$/){
                    if($2 eq "to"){
                        $size = "0|";
                    }else{
                       $size = "$3|$1";
                    }
                }else{
                    $size = $size."|";
                }

                print " |$state|$name|$size\n";
            }
            if ($event == 1 && $_ =~ m/\s+(.*?)\s+Subevent Code for/){
                print " |s|$1\n";
            }
            if ($_ =~ m/Condition for (.*?) = 0x(.*?)/){
                #Abort and handle manually
                print "a\n";
                $state = 0;
            }
            if ($_ =~ m/^Return parameters:/){
                $state = 4;
            }
            if ($_ =~ m/generated \(unless/){
                $state = 5;
            }
        }
        if($state == 5){
            if ($_ =~ m/HCI_Command_Complete/){
                print " |5|HCI_Command_Complete\n";
                $state = 0;
            }
            if ($_ =~ m/HCI_Command_Status/ || $_ =~ m/Command Status event/){
                print " |5|HCI_Command_Status\n";
                $state = 0;
            }
        }
    }
}

close FILE;
