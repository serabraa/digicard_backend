# DigiCard Backend

DigiCard is a Go backend that generates and signs Apple Wallet `.pkpass` files for digital business cards.

---

## Overview

The backend:

1. Receives a pass-generation request.
2. Creates the Apple Wallet pass data.
3. Loads the required images and assets.
4. Generates `manifest.json`.
5. Signs the manifest using the Apple Pass Type ID certificate.
6. Packages the files into a `.pkpass` archive.
7. Returns the signed pass to the client.

---

## API Endpoints

### Generate Pass

```http
POST /generate-pass
```

Generates and returns a signed Apple Wallet `.pkpass` file.

### Health Check

```http
GET /health
```

Successful response:

```text
ok
```

---

## Project Entry Point

The application starts from:

```text
cmd/api/main.go
```

The signing service is initialized here:

```go
signerService := signer.NewOpenSSLSigner(
    certPaths.passCertPath,
    certPaths.passKeyPath,
    certPaths.wwdrPath,
    filepath.Join(os.TempDir(), "digicard-work"),
)
```

The signer is then passed to the pass-generation service:

```go
generatePassService := generatepass.NewService(
    templateProvider,
    assetProvider,
    signerService,
    packagerService,
)
```

The detailed OpenSSL signing implementation is located in:

```text
infrastructure/signer/
```

---

## Pass Generation Flow

```text
POST /generate-pass
        │
        ▼
GeneratePassHandler
        │
        ▼
generatepass.Service
        │
        ├── Creates pass.json
        ├── Loads pass assets
        ├── Creates manifest.json
        ├── Signs the manifest
        └── Packages the files
        │
        ▼
Signed .pkpass response
```

---

## Required Environment Variables

The backend requires the following signing secrets:

```env
PASS_CERT_PEM=
PASS_KEY_PEM=
WWDR_PEM=
```

### `PASS_CERT_PEM`

The Apple **Pass Type ID certificate** in PEM format.

This is the main certificate used to sign DigiCard Apple Wallet passes.

### `PASS_KEY_PEM`

The private key associated with `PASS_CERT_PEM`.

The certificate and private key must belong to the same pair.

### `WWDR_PEM`

Apple’s Worldwide Developer Relations intermediate certificate.

It is used to create the certificate chain required for Apple Wallet signing.

---

## Temporary Certificate Files

During application startup, `prepareSignerFiles()` reads the certificate values from the environment and writes them to temporary files.

```text
/tmp/digicard-certs/pass-cert.pem
/tmp/digicard-certs/pass-key.pem
/tmp/digicard-certs/wwdr.pem
```

Permissions:

```text
Directory: 0700
Files:     0600
```

The application also converts escaped newline values such as `\n` into real line breaks before writing the PEM files.

---

# Apple Wallet Certificate Renewal

> [!IMPORTANT]
> The current DigiCard Pass Type ID certificate expires on **April 18, 2027**.

> [!WARNING]
> Begin the renewal process by **March 18, 2027**.  
> Do not wait until the certificate has already expired.

---

## Current Certificate Information

| Item | Value |
|---|---|
| Certificate type | Pass Type ID Certificate |
| Expiration date | **April 18, 2027** |
| Recommended renewal date | **March 18, 2027** |
| Certificate environment variable | `PASS_CERT_PEM` |
| Private-key environment variable | `PASS_KEY_PEM` |
| Intermediate certificate variable | `WWDR_PEM` |
| Certificate management | Apple Developer portal |
| Local certificate storage | macOS Keychain Access |

Update this table after every certificate renewal.

---

## What Happens When the Certificate Expires?

After the Pass Type ID certificate expires:

- New `.pkpass` files may fail to sign.
- Newly generated passes may not be accepted by Apple Wallet.
- Pass updates may stop working.
- The backend may return OpenSSL signing errors.
- Passes already installed on users’ devices should normally remain visible.

The Pass Type ID itself does not expire and should remain unchanged.

Example:

```text
pass.com.DigiCard.digital
```

---

## What Must Be Renewed?

### Pass Type ID Certificate

The main item that must be renewed is:

```env
PASS_CERT_PEM
```

Create the replacement certificate for the existing DigiCard Pass Type ID.

### Matching Private Key

The renewed certificate must have a matching private key:

```env
PASS_KEY_PEM
```

When a new Certificate Signing Request is generated, a new private key may also be created.

In that case, update both:

```env
PASS_CERT_PEM
PASS_KEY_PEM
```

### Apple Developer Program Membership

The Apple Developer Program membership must remain active.

Membership renewal and Pass Type ID certificate renewal are separate processes.

### WWDR Certificate

The WWDR certificate does not normally need annual renewal.

Only update:

```env
WWDR_PEM
```

when Apple replaces the required intermediate certificate or the current WWDR certificate expires.

---

# Certificate Renewal Procedure

## 1. Create a Certificate Signing Request

On macOS:

1. Open **Keychain Access**.
2. Open **Certificate Assistant**.
3. Select **Request a Certificate From a Certificate Authority**.
4. Enter the Apple Developer account email.
5. Select **Saved to disk**.
6. Save the generated `.certSigningRequest` file.

This process creates a private key in Keychain Access.

---

## 2. Create a New Pass Type ID Certificate

In the Apple Developer portal:

1. Open **Certificates, Identifiers & Profiles**.
2. Open **Certificates**.
3. Click the add button.
4. Select **Pass Type ID Certificate**.
5. Select the existing DigiCard Pass Type ID.
6. Upload the Certificate Signing Request.
7. Generate the certificate.
8. Download the `.cer` file.

> [!CAUTION]
> Do not create a new Pass Type ID unless the DigiCard identifier is intentionally changing.

---

## 3. Install the Certificate

Double-click the downloaded `.cer` file.

It should appear in:

```text
Keychain Access → My Certificates
```

Expand the certificate and confirm that a private key appears underneath it:

```text
Pass Type ID Certificate
└── Private Key
```

If the private key is missing, the certificate cannot be exported correctly for backend signing from that Mac.

---

## 4. Export the Certificate

In Keychain Access:

1. Select the Pass Type ID certificate.
2. Confirm that the private key is attached.
3. Export it as a `.p12` file.
4. Protect the export with a temporary password.
5. Store it securely.

Example filename:

```text
DigiCard-pass-certificate.p12
```

---

## 5. Convert the Certificate to PEM

Export the certificate:

```bash
openssl pkcs12 \
  -in DigiCard-pass-certificate.p12 \
  -clcerts \
  -nokeys \
  -out pass-cert.pem
```

Export the private key:

```bash
openssl pkcs12 \
  -in DigiCard-pass-certificate.p12 \
  -nocerts \
  -nodes \
  -out pass-key.pem
```

---

## 6. Check the New Certificate

Display certificate information:

```bash
openssl x509 \
  -in pass-cert.pem \
  -noout \
  -subject \
  -issuer \
  -serial \
  -dates
```

Display only the expiration date:

```bash
openssl x509 \
  -in pass-cert.pem \
  -noout \
  -enddate
```

Example output:

```text
notAfter=Apr 18 23:59:59 2027 GMT
```

---

## 7. Verify That the Certificate and Key Match

Generate the certificate public-key hash:

```bash
openssl x509 \
  -in pass-cert.pem \
  -pubkey \
  -noout |
openssl sha256
```

Generate the private-key public-key hash:

```bash
openssl pkey \
  -in pass-key.pem \
  -pubout |
openssl sha256
```

The two SHA-256 values must be identical.

---

## 8. Update Backend Secrets

Replace the deployed values for:

```env
PASS_CERT_PEM
PASS_KEY_PEM
```

Update `WWDR_PEM` only when necessary.

The complete PEM values must be stored, including their headers and footers.

Certificate format:

```text
-----BEGIN CERTIFICATE-----
...
-----END CERTIFICATE-----
```

Private-key format:

```text
-----BEGIN PRIVATE KEY-----
...
-----END PRIVATE KEY-----
```

---

## 9. Redeploy the Backend

Restart or redeploy the backend after updating the secrets.

The signing secrets are read during application startup, so updating the secret values without restarting the application may not update the running instance.

---

## 10. Test the Renewed Certificate

Check the health endpoint:

```http
GET /health
```

Generate a new pass:

```http
POST /generate-pass
```

Confirm that:

- The backend returns a `.pkpass` file.
- No OpenSSL signing errors occur.
- The pass can be added to a real iPhone.
- The Pass Type ID is correct.
- The organization name is correct.
- Pass images and fields appear correctly.
- Pass updates still work, when supported.

---

# Annual Renewal Checklist

Complete this checklist before **April 18, 2027**.

- [ ] Confirm that the Apple Developer membership is active.
- [ ] Check the Pass Type ID certificate expiration date.
- [ ] Confirm which certificate is deployed in production.
- [ ] Create a new Certificate Signing Request.
- [ ] Generate a certificate for the existing Pass Type ID.
- [ ] Download and install the new certificate.
- [ ] Confirm that the private key is visible in Keychain Access.
- [ ] Export the certificate and key as `.p12`.
- [ ] Convert the certificate to `pass-cert.pem`.
- [ ] Convert the private key to `pass-key.pem`.
- [ ] Verify that the certificate and private key match.
- [ ] Update `PASS_CERT_PEM`.
- [ ] Update `PASS_KEY_PEM`.
- [ ] Check whether `WWDR_PEM` needs to be updated.
- [ ] Redeploy the backend.
- [ ] Generate a test `.pkpass`.
- [ ] Add the test pass to a real iPhone.
- [ ] Confirm the new expiration date.
- [ ] Update this README with the next renewal date.
- [ ] Create a reminder one month before the next expiration date.

---

# Security

> [!CAUTION]
> Never commit certificates, private keys, environment files, or `.p12` files to Git.

Recommended `.gitignore` entries:

```gitignore
.env
.env.*
*.p12
*.pem
digicard-certs/
digicard-work/
```

Sensitive files include:

```text
pass-cert.pem
pass-key.pem
wwdr.pem
DigiCard-pass-certificate.p12
.env
```

Store production secrets in a secure secret-management system.

Anyone who has access to both the Pass Type ID certificate and its private key may be able to sign passes using your Pass Type ID.

---

# Emergency Recovery

If the certificate has already expired:

1. Keep the existing Pass Type ID.
2. Generate a new Certificate Signing Request.
3. Create a new Pass Type ID certificate.
4. Install the certificate in Keychain Access.
5. Confirm that its private key is attached.
6. Export the certificate and private key.
7. Convert them to PEM.
8. Update `PASS_CERT_PEM`.
9. Update `PASS_KEY_PEM`.
10. Redeploy the backend.
11. Generate and test a new `.pkpass`.

If the original private key has been lost, create a new private key and Certificate Signing Request.

The original private key is not required for generating new passes after the backend is updated with a new valid certificate and its matching private key.

---

# Renewal History

Keep a record of certificate changes here.

| Certificate | Renewal date | Expiration date | Updated by | Notes |
|---|---|---|---|---|
| Current certificate | — | April 18, 2027 | — | Initial certificate record |
| Next certificate | — | — | — | Add after renewal |

---

## Reminder

```text
Renew DigiCard Pass Type ID certificate by March 18, 2027.
Current certificate expires on April 18, 2027.
```
