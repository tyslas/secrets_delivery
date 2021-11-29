#!/bin/bash


errArr=()

if [ -z "$1" ]; then
  errArr+=("must submit location of private key as the first argument")
else
  [ -f "$1" ] && echo "using $1 as private key" || errArr+=("$1 does not exist")
fi

if [ -z "$2" ]; then
  errArr+=("must submit location of public key as the second argument")
else
  [ -f "$2" ] && echo "using $2 as public key" || errArr+=("$2 does not exist")
fi

if [ -z "$3" ]; then
  errArr+=("must submit location of message to encrypt or decrypt as the third argument")
else
  [ -f "$3" ] && echo "using $3 as message for encryption/decryption" || errArr+=("$3 does not exist")
fi

if [ -z "$4" ]; then
  errArr+=("must submit \"enc\" or \"dec\" to indicate encryption or decryption as the fourth argument")
else
  if [ "$4" != "enc" ] && [ "$4" != "dec" ]; then
    errArr+=("only valid values for the fourth argument are \"enc\" or \"dec\"")
  fi
  if [ "$4" == "dec" ]; then
    if [ -z "$5" ]; then
      errArr+=("must submit location of the signature as the fifth argument")
    else
      [ -f "$5" ] && echo "using $5 as signature to verify" || errArr+=("$5 does not exist")
    fi
  fi
fi

if [ ${#errArr[@]} -gt 0 ]; then
  echo "errors:"
  printf '    %s\n' "${errArr[@]}"
  exit 1
fi

if [ "$4" == "enc" ]; then
  openssl rsautl -encrypt -inkey $2 -pubin -in $3 -out ./payload/encrypt.dat
  openssl dgst -sha256 -sign $1 -out ./payload/sign.sha256 ./payload/encrypt.dat
  openssl base64 -in ./payload/sign.sha256 -out ./payload/signature.dat
  tar czf ./payload/payload.tgz ./payload
  if [ -f ./payload/payload.tgz ]; then
    scp -i ~/.ssh/tito_aws.pem ./payload/payload.tgz ec2-user@ec2-54-193-103-167.us-west-1.compute.amazonaws.com:/home/ec2-user
  else
    echo "    failed to create tar for payload"
    exit 1
  fi
fi

if [ "$4" == "dec" ]; then
  echo "verifying signature and decrypting"
  tar xvzf payload.tgz
  openssl base64 -d -in $5 -out ./payload/sign.sha256
  openssl dgst -sha256 -verify $2 -signature ./payload/sign.sha256 ./payload/encrypt.dat
  status=$?
  if [ $status -eq 0 ]; then
    echo "signature verified... decrypting"
    openssl rsautl -decrypt -inkey $1 -in ./payload/encrypt.dat -out ./decrypt.txt 
    cat ./decrypt.txt
  else
    echo "    unable to verify signature"
    exit 1
  fi
fi


