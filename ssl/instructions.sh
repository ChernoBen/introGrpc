#!/bin/bash
#mude este CN para o host que utiliza em seu ambiente
#SERVER_CN=localhost

#1: Gerar certificado de autoridade
openssl genrsa -passout pass:1111 -des3 -out ca.key 4096
openssl req -passin pass:1111 -new -x509 -days 365 -key ca.key -out ca.crt -subj "/CN=${SERVER_CN}"

#2: gerar server private key
openssl genrsa -passout pass:1111 -des3 -out server.key 4896

#3: obter certificado assinando request do CA (server.csr)
openssl req -passin pass:1111 -new -key server.key -out server.csr -subj "/CN=${SERVER_CN}"

#4: assinar certificado com CA
openssl x509 -req -passin pass:1111 -days 365 -in server.csr -CA ca.crt -CAkey ca.key -set_serial 01 -out server.crt

#5: converter o certificado do server para .pem formt (server)
openssl pkcs8 -topk8 -nocrypt -passin pass:1111 -in server.key -out server.pem