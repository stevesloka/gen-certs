# gen-certs

The goal of this repo is to provide a way to generate certs

# Example

Gen CA:
```
$ cfssl genkey -initca config/ca-csr.json | cfssljson -bare certs/ca  
```

Gen Node Cert:
```
$ cfssl gencert \
	  -ca certs/ca.pem \
	  -ca-key certs/ca-key.pem \
	  -config config/ca-config.json \
	  config/req-csr.json | cfssljson -bare certs/node
```

Convert CA to pkcs12:
```
$ openssl pkcs12 -export -inkey certs/ca-key.pem -in certs/ca.pem -out certs/ca.pkcs12 -password pass:changeit
```

Convert Node to pkcs12:
```
$ openssl pkcs12 -export -inkey certs/node-key.pem -in certs/node.pem -out certs/node.pkcs12 -password pass:changeit
```

Create JKS:
```
$ keytool -importkeystore -srckeystore certs/ca.pkcs12 -srcalias '1' -destkeystore certs/truststore.jks \
   -storepass "changeit" -srcstoretype pkcs12 \
   -srcstorepass "changeit" -destalias elasticsearch-ca

$ keytool -importkeystore -srckeystore certs/node.pkcs12 -srcalias '1' -destkeystore certs/node-keystore.jks \
   -storepass "changeit" -srcstoretype pkcs12 \
   -srcstorepass "changeit" -destalias elasticsearch-node
```

# Validate 

```
$ openssl x509 -noout -text -in certs/node.pem    
```