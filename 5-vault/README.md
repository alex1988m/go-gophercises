# Vault - Encrypted Key-Value Store

Vault is a simple command-line application that provides an encrypted key-value store. It allows users to securely store and retrieve sensitive information using AES encryption.

## Features

- Set encrypted values associated with keys
- Retrieve decrypted values using keys
- AES encryption for secure storage
- File-based persistence
- Command-line interface

## Installation

1. Clone the repository:
   ```
   git clone https://github.com/alex1988m/go-gophercises.git
   ```

2. Navigate to the project directory:
   ```
   cd go-gophercises/5-vault
   ```

3. Install dependencies:
   ```
   go mod tidy
   ```

4. Create a `.env` file in the project root and set your cipher key:
   ```
   CIPHER_KEY=your32charactercipherkeygoeshere
   ```

## Usage

### Set a value

```
go run main.go set <key> <value>
```

Example:
```
go run main.go set mypassword secretpassword123
```

### Get a value

```
go run main.go get <key>
```

Example:
```
go run main.go get mypassword
```

## Project Structure

- `main.go`: Entry point of the application, handles CLI commands
- `vault/vault.go`: Core vault functionality including encryption/decryption
- `vault/file_vault.go`: File-based storage implementation
- `logger/logger.go`: Logging configuration

## Dependencies

- [github.com/urfave/cli/v2](https://github.com/urfave/cli): CLI application framework
- [github.com/sirupsen/logrus](https://github.com/sirupsen/logrus): Structured logger
- [github.com/joho/godotenv](https://github.com/joho/godotenv): Loads environment variables from .env file
- [github.com/pkg/errors](https://github.com/pkg/errors): Error handling

## Security Note

Ensure that your `.env` file and `vault.json` are not shared or committed to version control, as they contain sensitive information.
