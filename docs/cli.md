# CLI


```sh
# Login with AppRole
./cryptkeeper-cli login auth/approle --role_id=<Role ID> --secret_id=<Secret ID>

# Create a new secret
./cryptkeeper-cli create secret --path=/kvstore --key=mysecret --value="supersecretvalue"
```