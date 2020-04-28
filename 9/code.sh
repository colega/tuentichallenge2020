#!/usr/bin/env bash

    chr() {
        printf \\$(printf '%03o' $1)
    }

    function hex() {
        printf '%02X\n' $1
    }

    function encrypt() {
        key=$1
        msg=$2
        crpt_msg=""
        for ((i=0; i<${#msg}; i++)); do
            c=${msg:$i:1}
            asc_chr=$(echo -ne "$c" | od -An -tuC)
            key_pos=$((${#key} - 1 - ${i}))
            key_char=${key:$key_pos:1}
            crpt_chr=$(( $asc_chr ^ ${key_char} ))
            hx_crpt_chr=$(hex $crpt_chr)
            crpt_msg=${crpt_msg}${hx_crpt_chr}
            echo "c=${c}, asc_chr=${asc_chr}, key_pos=${key_pos}, key_char=${key_char}, crpt_char=${crpt_chr}, hex=${hx_crpt_chr}"
        done
        echo $crpt_msg
    }
