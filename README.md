# budgetto

bash```
openssl genpkey -algorithm RSA -out private_key.pem -pkeyopt rsa_keygen_bits:2048
openssl rsa -pubout -in private_key.pem -out public_key.pem

base64 -w 0 private_key.pem > private_key_base64.txt
```
