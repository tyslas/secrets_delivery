#!/bin/bash

openssl genrsa -out ~/.ssh/my_private_key.pem 4096
openssl rsa -in ~/.ssh/my_private_key.pem -out ~/.ssh/my_public_key.pem -outform PEM -pubout