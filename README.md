# env-vault

Launch a subprocess with environment variables from an encrypted file

## Usage

Unix "standard" allows appending a double dash "--" to the command
and everything after that can get passed to a subcommand.

For example:

```sh
env-vault --vault prod.env docker-compose -- run --rm db env
```

This will decode secrents in `prod.env` and expose them to
`docker-compose`. The command otherwise would look like so:

```
docker-compose run --rm db env
```
