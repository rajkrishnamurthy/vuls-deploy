#!/bin/sh
trap ctrl_c INT

function ctrl_c() {}

if [ ! -f /vuls/data/oval.sqlite3 ]; then 
    mkdir data; > oval.sqlite3;goval-dictionary fetch-ubuntu -dbpath=/vuls/data/oval.sqlite3 18; 
fi
if [ ! -f /vuls/data/cve.sqlite3 ]; then 
    > cve.sqlite3;
    for i in `seq 2002 $(date +"%Y")`; 
    do 
        go-cve-dictionary fetchnvd -dbpath /vuls/data/cve.sqlite3 -years $i; 
    done;
    goval-dictionary fetch-debian -dbpath=/vuls/data/oval.sqlite3 7 8 9 10; 
    goval-dictionary fetch-redhat -dbpath=/vuls/data/oval.sqlite3 5 6 7 8;
    goval-dictionary fetch-ubuntu -dbpath=/vuls/data/oval.sqlite3 14 16 18 19;
    goval-dictionary fetch-suse -dbpath=/vuls/data/oval.sqlite3 -opensuse 13.2;
    goval-dictionary fetch-oracle -dbpath=/vuls/data/oval.sqlite3;
    goval-dictionary fetch-alpine -dbpath=/vuls/data/oval.sqlite3 3.3 3.4 3.5 3.6;
    goval-dictionary fetch-amazon -dbpath=/vuls/data/oval.sqlite3; 
fi

vuls server -listen=localhost:80
