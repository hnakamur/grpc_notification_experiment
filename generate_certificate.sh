#!/bin/sh -x
set -eu
ssl_ca_dir=ssl/ca
ssl_server_dir=ssl/server

mkdir -p "$ssl_ca_dir"/certs
mkdir -m 700 -p "$ssl_ca_dir"/private
mkdir -p "$ssl_ca_dir"/crl
mkdir -p "$ssl_ca_dir"/newcerts
if [ ! -f "$ssl_ca_dir"/serial ]; then
  echo '01' > "$ssl_ca_dir"/serial
fi
if [ ! -f "$ssl_ca_dir"/index.txt ]; then
  touch "$ssl_ca_dir"/index.txt
fi
cat <<EOF > "$ssl_ca_dir"/config
[ ca ]
default_ca      = CA_default            # The default ca section

[ CA_default ]
dir = $ssl_ca_dir
database = \$dir/index.txt
new_certs_dir = \$dir/newcerts

certificate    = \$dir/cacert.pem
serial         = \$dir/serial
private_key    = \$dir/private/cakey.pem

RANDFILE       = \$dir/private/.rand    # random number file

default_days   = 365                   # how long to certify for
default_crl_days= 30                   # how long before next CRL
defaults_bits  = 2048
default_md     = sha256

policy         = policy_any            # default policy
email_in_dn    = no                    # Don't add the email into cert DN

name_opt       = ca_default            # Subject name display option
cert_opt       = ca_default            # Certificate display option
copy_extensions = none                 # Don't copy extensions from request

[ policy_any ]
countryName            = supplied
stateOrProvinceName    = optional
organizationName       = optional
organizationalUnitName = optional
commonName             = supplied
emailAddress           = optional
EOF

if [ ! -f "$ssl_ca_dir"/private/cakey.pem ]; then
  (cd "$ssl_ca_dir" && openssl req -new -x509 -newkey rsa:2048 -out cacert.pem -keyout private/cakey.pem -days 365)
  chmod 600 "$ssl_ca_dir"/private/cakey.pem
fi

mkdir -p "$ssl_server_dir"
if [ ! -f "$ssl_server_dir"/server.crt ]; then
  if [ ! -f "$ssl_server_dir"/server.key ]; then
    openssl req -new -keyout "$ssl_server_dir"/server.key -out "$ssl_server_dir"/server.csr
    openssl rsa -in "$ssl_server_dir"/server.key -out "$ssl_server_dir"/server.key
  fi
  openssl ca -config "$ssl_ca_dir"/config -out "$ssl_server_dir"/server.crt -infiles "$ssl_server_dir"/server.csr
fi
