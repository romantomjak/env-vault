# env-vault

env-vault provides a convenient way to launch a program
with environment variables populated from an encrypted file.
Inspired by [ansible-vault](https://docs.ansible.com/ansible/latest/cli/ansible-vault.html)
and [envconsul](https://github.com/hashicorp/envconsul).

## Installation

Download and install using go get:

```sh
go get -u github.com/romantomjak/env-vault
```

or grab a binary from [releases](https://github.com/romantomjak/env-vault/releases/latest) section!

## Usage

env-vault allows to do pretty cool stuff like injecting secrets into docker-compose files, but in general use it like so:

```
env-vault <vault> <program>
```

The `<program>` argument is the executable that will be launched with environment variables from the encrypted file pointed to by `<vault>` argument.

env-vault takes advantage of a POSIX standard that uses `--` to signify the end of command line options. Meaning that everything after that can get passed on to a sub-command. Use this one one weird trick to pass arguments to programs:

```
env-vault <vault> <program> -- <program-arg1> <program-arg2> <...>
```

## Example Docker Compose Use

This section describes how to use env-vault to safely store secrets and then pass them on to docker-compose.

### docker-compose.yml

Please note the `POSTGRES_PASSWORD` field that is set only as a key. This tells docker-compose to resolve the value on the machine on which the compose is running on.

```yml
services:
  db:
    image: postgres:13-alpine
    environment:
      - POSTGRES_USER=myproject
      - POSTGRES_PASSWORD
```

### prod.env

This is an encrypted file created with env-vault. Vaults are created with the `create` sub-command and can be viewed using the `view` sub-command. Let's create a vault named `prod.env` to hold our production secrets:

```sh
env-vault create prod.env
```

Running the command above will prompt for a new password and then open your favorite `$EDITOR` to input the environment variables. Let's add the `POSTGRES_PASSWORD` secret now:

```ini
POSTGRES_PASSWORD=passwordformyproject
```

Save the file and close your editor. env-vault will encrypt the plain text using AES256-GCM symmetrical encryption.

### docker-compose

env-vault takes advantage of a POSIX standard that uses `--` to signify the end of command line options. Meaning that everything after that can get passed on to a sub-command. Let's see how we can use that to decrypt secrets for docker-compose:

```sh
env-vault prod.env docker-compose -- up -d
```

It looks somewhat mad, but essentially env-vault will decrypt `prod.env` and expose found environment variables to
`docker-compose`. The command would otherwise look like so:

```
docker-compose up -d
```

Now lets inspect the container to see that the `POSTGRES_PASSWORD` environment variable was successfully injected:

```sh
$ docker-compose exec db env | grep POSTGRES_PASSWORD
POSTGRES_PASSWORD=passwordformyproject
```

Wooo!!! :rocket:

## Advanced Use

You can further ease the workflow by setting the `ENV_VAULT_PASSWORD` environment variable. env-vault will use it to decrypt vaults automatically and you won't get prompted for a password any time you interact with env-vault.

```sh
export ENV_VAULT_PASSWORD=somepassword
```

This works exceptionally well when docker-compose is used with Makefiles:

```Makefile
up:
	env-vault prod.env docker-compose -- up -d
```

so now you can run:

```sh
make up
```

and env-vault will take it from there to make it work automagically! :sparkles:

## Contributing

You can contribute in many ways and not just by changing the code! If you have
any ideas, just open an issue and tell me what you think.

Contributing code-wise - please fork the repository and submit a pull request.

## License

MIT