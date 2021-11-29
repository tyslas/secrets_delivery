#!/bin/bash


errArr=()

if [ -z "$1" ]; then
  errArr+=("must submit location of private key as the first argument")
else
  [ -f "$1" ] && echo "using $1 as private key" || errArr+=("$1 does not exist")
  PR_KEY=$1
fi

if [ -z "$2" ]; then
  errArr+=("must submit location of public key as the second argument")
else
  [ -f "$2" ] && echo "using $2 as public key" || errArr+=("$2 does not exist")
  PB_KEY=$2
fi

if [ -z "$3" ]; then
  errArr+=("must submit location of message to encrypt or decrypt as the third argument")
else
  echo "using $3 as message for encryption/decryption"
  MSG=$3
fi

if [ -z "$4" ]; then
  errArr+=("must submit \"enc\" or \"dec\" to indicate encryption or decryption as the fourth argument")
else
  if [ "$4" != "enc" ] && [ "$4" != "dec" ]; then
    errArr+=("only valid values for the fourth argument are \"enc\" or \"dec\"")
  fi
  if [ "$4" == "enc" ]; then
    if [ -z "$5" ]; then
      errArr+=("must submit username on remote server for fifth argument")
    else
      echo "using $5 as username for remote server"
      USER=$5
    fi
    if [ -z "$6" ]; then
      errArr+=("must submit SSH private key to scp files onto remote server for sixth argument")
    else
      echo "using $6 as SSH private key for remote server"
      echo "    NOTE: this can be the same as the private key used for encryption"
      SSH_KEY=$6
    fi
    if [ -z "$7" ]; then
      errArr+=("must submit IP of remote server for seventh argument")
    else
      echo "using $7 as IP for remote server"
      IP=$7
    fi
    if [ -z "$8" ]; then
      errArr+=("must submit path to copy data to on remote server for eighth argument")
    else
      echo "using $8 as IP for remote server"
      DEST_DIR=$8
    fi
  fi
  if [ "$4" == "dec" ]; then
    if [ -z "$5" ]; then
      errArr+=("must submit location of the signature as the fifth argument")
    else
      echo "using $5 as signature to verify"
      SIG_LOC=$5
    fi
  fi
fi

if [ ${#errArr[@]} -gt 0 ]; then
  echo "errors:"
  printf '    %s\n' "${errArr[@]}"
  exit 1
fi

if [ "$4" == "enc" ]; then
  openssl rsautl -encrypt -inkey "$PB_KEY" -pubin -in "$MSG" -out ./payload/encrypt.dat
  openssl dgst -sha256 -sign "$PR_KEY" -out ./payload/sign.sha256 ./payload/encrypt.dat
  openssl base64 -in ./payload/sign.sha256 -out ./payload/signature.dat
  tar czf ./payload/payload.tgz ./payload
  if [ -f ./payload/payload.tgz ]; then
    scp -i "$SSH_KEY" ./payload/payload.tgz "$USER"@"$IP":"$DEST_DIR"
  else
    echo "    failed to create tar for payload"
    exit 1
  fi
fi

if [ "$4" == "dec" ]; then
  echo "verifying signature and decrypting"
  tar xvzf payload.tgz
  openssl base64 -d -in "$SIG_LOC" -out ./payload/sign.sha256
  openssl dgst -sha256 -verify "$PB_KEY" -signature ./payload/sign.sha256 ./payload/encrypt.dat
  status=$?
  if [ $status -eq 0 ]; then
    echo "signature verified... decrypting"
    openssl rsautl -decrypt -inkey "$PR_KEY" -in ./payload/encrypt.dat -out ./decrypt.txt
    cat ./decrypt.txt
  else
    echo "    unable to verify signature"
    exit 1
  fi
fi


