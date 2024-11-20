```sh
./cryptkeeper-cli
A CLI for managing secrets with CryptKeeper

Usage:
  cryptkeeper-cli [command]

Available Commands:
  assign-access Assign access to a secret
  completion    Generate the autocompletion script for the specified shell
  create        Create a new secret
  detect        Detect secrets in a specified directory
  get           Retrieve a secret
  help          Help about any command
  import        Import secrets from config files
  login         Authenticate using AppRole
  paths         Retrieve Paths
  rotate        Rotate a secret
  secrets       Retrieve Secrets for a given path
```



go build -o cryptkeeper-cli


Run the CLI commands:

**Login**
```sh
./cryptkeeper-cli login --role_id=7bf3aeb6-d9c2-410c-b23e-0b8f8e85dc4a --secret_id=a77eb0f8-ed18-45b8-800a-7f7b1cee8a86
export TOKEN="your_jwt_token"
export CRYPTKEEPER_TOKEN="your_jwt_token"
```

**Create a Secret**

```sh
./cryptkeeper-cli create create /kvstore foo bar
./cryptkeeper-cli create /kvstore foo bar --expires-at "2024-05-28T23:59:59Z"

```

**Rotate a Secret**
```sh
./cryptkeeper-cli rotate /path/to/secret "my_new_secret_value"
```


**Get a Secret**
```sh
./cryptkeeper-cli get /dev /attribute/test 
./cryptkeeper-cli get /dev /attribute/test --version=1
```
