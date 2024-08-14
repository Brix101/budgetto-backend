# budgetto

## Generating RSA Keys

To generate RSA keys and convert them to base64 format, follow these steps:

1. **Generate a private key**:

   ```bash
   openssl genpkey -algorithm RSA -out private_key.pem -pkeyopt rsa_keygen_bits:2048
   ```

2. **Generate the corresponding public key**:

   ```bash
   openssl rsa -pubout -in private_key.pem -out public_key.pem
   ```

3. **Convert the private key to base64 format**:

   ```bash
   base64 -w 0 private_key.pem > private_key_base64.txt
   ```

4. **Convert the public key to base64 format**:
   ```bash
   base64 -w 0 public_key.pem > public_key_base64.txt
   ```

These commands will create `private_key.pem`, `public_key.pem`, `private_key_base64.txt`, and `public_key_base64.txt` files in your current directory.
