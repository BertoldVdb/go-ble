#!/usr/bin/perl

use strict;

my $inStruct;
my $outStruct;
my $outDecode;
my $inEncode;
my $name;
my $hasInput = 0;
my $hasOutput = 0;
my $hasStatus = 0;
my $eventStatus = 0;
my $ogf = "";
my $ocf = "";
my $point = "";
my $varcnt = 0;
my $desired = int(@ARGV[1]);

print <<EOF;
package hcicommands

import (
\t"encoding/binary"
\thcicmdmgr "github.com/BertoldVdb/go-ble/hci/cmdmgr"
)

EOF


sub getTopic{
    my $id=int(@_[0]);
    if ($id == 1) { return "LinkControl"; }
    if ($id == 2) { return "LinkPolicy"; }
    if ($id == 3) { return "Baseband"; }
    if ($id == 4) { return "Informational"; }
    if ($id == 5) { return "Status"; }
    if ($id == 6) { return "Testing"; }
    if ($id == 7) { return ""; }
    if ($id == 8) { return ""; } #The command already has a prefix 
    return "Unknown"
}

sub clean{
    my $name = @_[0];
    if ($name =~ /^\d/){
        $name = "I$name";
    }
    $name =~ tr/ //ds;
    $name =~ tr/\_//ds;
    $name =~ tr/-//ds;
    return $name;
}

sub lenToType {
    my $len=int(@_[0]);
    if ($len == 0) { return "[]byte"; }
    if ($len == 1) { return "uint8"; }
    if ($len == 2) { return "uint16"; }
    if ($len == 3) { return "uint32"; }
    if ($len == 4) { return "uint32"; }
    if ($len == 8) { return "uint64"; }
    return "[$len]byte";
}

sub decodeToType {
    my $len=int(@_[0]);
    my $fn = @_[1];
    my $lenVar = @_[2];
    if ($len == 0) { 
        if ($lenVar eq "" ){
            return "$fn = append($fn"."[:0], r.GetRemainder()...)"; 
        }else{
            return "$fn = append($fn"."[:0], r.Get($lenVar)...)"; 
        }
    }
    if ($len == 1) { return "$fn = r.GetOne()"; }
    if ($len == 2) { return "$fn = binary.LittleEndian.Uint16(r.Get(2))"; }
    if ($len == 3) { return "$fn = decodeUint24(r.Get(3))"; }
    if ($len == 4) { return "$fn = binary.LittleEndian.Uint32(r.Get(4))"; }
    if ($len == 8) { return "$fn = binary.LittleEndian.Uint64(r.Get(8))"; }
    return "copy($fn"."[:], r.Get($len))";
}

sub encodeToType {
    my $len=int(@_[0]);
    my $fn = @_[1];
    if ($len == 0) { return "w.PutSlice($fn)"; }
    if ($len == 1) { return "w.PutOne($fn)"; }
    if ($len == 2) { return "binary.LittleEndian.PutUint16(w.Put(2), $fn)"; }
    if ($len == 3) { return "encodeUint24(w.Put(3), $fn)"; }
    if ($len == 4) { return "binary.LittleEndian.PutUint32(w.Put(4), $fn)"; }
    if ($len == 8) { return "binary.LittleEndian.PutUint64(w.Put(8), $fn)"; }
    return "copy(w.Put($len), $fn"."[:])";
}

sub finalizeStruct {
    my $newName = @_[0];

    my $comment = "Section 7.$ogf.$point";

    my $out = "Output";
    if($ogf == 7){
        $out = "Event";
        $outStruct = "// $newName"."$out represents the event specified in $comment\ntype $newName"."$out struct {\n$outStruct";
    }else{
        $outStruct = "// $newName"."$out represents the output of the command specified in $comment\ntype $newName"."$out struct {\n$outStruct";
    }

    $inStruct = "// $newName"."Input represents the input of the command specified in $comment\ntype $newName"."Input struct {\n$inStruct";
    $inEncode = "func (i $newName"."Input) encode(data []byte) []byte {\n\tw := writer{data: data};\n$inEncode";
    $outDecode = "func (o *$newName"."$out) decode(data []byte) bool {\n\tr := reader{data: data};\n$outDecode";
    
    my $params = "";
    my $return = "error";
    my $returnVal = "err";
   
    if ($hasInput) {
        print $inStruct."}\n\n";
        print $inEncode."\treturn w.Data()\n}\n\n";
        $params .= "params $name"."Input";
    }
    if ($hasOutput) {
        print $outStruct."}\n\n";
        print $outDecode."\treturn r.Valid()\n}\n\n";
        if(length($params)){
            $params .=", ";
        }
        $params .= "result *$name"."Output";
        $return = "(*$name"."Output, error)";
        $returnVal = "result, err";
    }

    if($ogf != 7){
        print "// $name"."Sync executes the command specified in $comment synchronously\n";
        print "func (c *Commands) $name"."Sync ($params) $return {\n";

        if($hasOutput){
            print <<EOF;
\tif result == nil {
\t\tresult = &$name\Output{}
\t}

EOF
        }

        print <<EOF;
\tbuffer, err := c.hcicmdmgr.CommandRunGetBuffer(0, hcicmdmgr.HCICommand{OGF: $ogf, OCF: $ocf}, nil)
\tif err != nil {
\t\treturn $returnVal
\t}

EOF
	
        if($hasInput){
            print "\tbuffer.Buffer = params.encode(buffer.Buffer)\n";
        }

        my $resp = "_, err =";
        if ($hasOutput || $hasStatus){
            $resp = "response, err :=";
        }

        print <<EOF;
\t$resp c.hcicmdmgr.CommandRunPutBuffer(buffer)
\tif err != nil {
\t\treturn $returnVal
\t}

EOF

        if($hasOutput){
            print <<EOF;
\tif !result.decode(response) {
\t\terr = ErrorMalformed
\t}

EOF
        }

        if($hasStatus){
            print "\terr = HciErrorToGo(response, err)\n";
        }

        print <<EOF;

\terr2 := c.hcicmdmgr.CommandRunReleaseBuffer(buffer)
\tif err2 != nil {
\t\terr = err2
\t}

\treturn $returnVal
}

EOF
    }
}

my %codeValues;
sub resetStruct(){
    $inStruct = "";
    $outStruct = "";
    $inEncode = "";
    $outDecode = "";
    $hasInput = 0;
    $hasOutput = 0;
    $hasStatus = 0;
    %codeValues = ();
}

sub codeValue(){
    my $sn = @_[0];
    my $part = @_[1];

    my $extra = "";
    my $ppart = $sn.$part;
    if (exists($codeValues{$ppart})) {
        return $extra,$codeValues{$ppart};
    }

    my $value = "int($sn.$part)";
    if($part =~ m/^Bitssetin(.*)$/){
        $extra = "\tvar$varcnt := countSetBits(uint64($sn.$1))\n";
        $value = "var$varcnt";
        $varcnt++;
    }elsif($part =~ m/^SUM\((.*)\[(.*)\]\)$/){
        $extra = "\tvar$varcnt := 0\n";
        $extra .= "\tfor _, m := range $sn.$1 {\n";
        $extra .= "\t\tvar$varcnt += int(m)\n";
        $extra .= "\t}\n";
        $value = "var$varcnt";
        $varcnt++;
    }
    $codeValues{$ppart} = $value;

    return $extra, $value
}

my $active = 0;

open (FILE, @ARGV[0]) or die "File not found";
while(<FILE>){
    chomp $_;
    chomp $_;

    my @parts = split /\|/, $_;
    if (@parts[0] eq "c"){
        &resetStruct();
        $name = &clean(@parts[3]);
        $name = &getTopic(@parts[1]) . $name;
        $ogf = @parts[1];
        $point = @parts[2];
        $ocf = @parts[4];
        $active = 0;
        if (int($ogf) == $desired || $desired<=0) {
            $active = 1;
        }
    }elsif($active == 1){
        my $array="";
        if (@parts[4] ne ""){
            $array .= "[]";
            @parts[2] = substr(@parts[2], 0, length(@parts[2])-3);
        }
        
        @parts[2] = &clean(@parts[2]);
        @parts[4] = &clean(@parts[4]);
        my $type = &lenToType(@parts[3]);

        my $value = "\t".@parts[2]." $array$type\n";
        if (@parts[1] eq "3"){
            $inStruct .= $value;
            $hasInput = 1;

            if (@parts[4] ne ""){
                (my $extra, my $value) = &codeValue("i", @parts[4]);
                my $test = "len(i.@parts[2]) != $value";

                $inEncode .= "$extra\tif $test {\n\t\tpanic(\"$test\")\n\t}\n";
                $inEncode .= "\tfor _, m := range i.@parts[2] {\n";
                $inEncode .= "\t\t".&encodeToType(@parts[3], "m")."\n";
                $inEncode .= "\t}\n";

            }else{
                $inEncode .= "\t".&encodeToType(@parts[3], "i.".@parts[2])."\n";
            }
        }
        if (@parts[1] eq "4"){
            $outStruct .= $value;
            if (@parts[2] ne "Status"){
                $hasOutput = 1;
            }else{
                $hasStatus = 1;
            }

            if (@parts[4] ne ""){
                (my $extra, my $value) = &codeValue("o", @parts[4]);
                $outDecode .= $extra;
                $outDecode .= "\tif cap(o.".@parts[2].") < $value {\n";
                $outDecode .= "\t\to.".@parts[2]." = make([]$type, 0, $value)\n";
                $outDecode .= "\t}\n";
                $outDecode .= "\to.".@parts[2]." = o.".@parts[2]."[:$value]\n";
                $outDecode .= "\tfor j:=0; j<$value; j++ {\n";
                $outDecode .= "\t\t".&decodeToType(@parts[3], "o.".@parts[2]."[j]", "int(o.DataLength[j])")."\n";
                $outDecode .= "\t}\n";

            }else{
                $outDecode .= "\t".&decodeToType(@parts[3], "o.".@parts[2], "")."\n";
            }
        }
        if (@parts[1] eq "5"){
            $eventStatus = "false";
            if (@parts[2] eq "HCICommandStatus") {
                $eventStatus = "true";
                $hasStatus = 1;
            }
            &finalizeStruct($name);
        }
    }
}
close (FILE);
