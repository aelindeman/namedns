# namedns

[![Build Status](https://travis-ci.org/aelindeman/namedns.svg?branch=master)](https://travis-ci.org/aelindeman/namedns)
[![Maintainability](https://api.codeclimate.com/v1/badges/02357c367b6f044ad810/maintainability)](https://codeclimate.com/github/aelindeman/namedns/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/02357c367b6f044ad810/test_coverage)](https://codeclimate.com/github/aelindeman/namedns/test_coverage)

A command-line utility to manipulate Name.com DNS records

## Install

Grab the latest binary from the [releases page](https://github.com/aelindeman/namedns/releases/latest) and drop it in your `$PATH`.

## Usage

A username and API key to the Name.com API are required.

### Getting an API key

You can get an API key by going to <https://name.com/reseller/apply> and filling out the form. Support will later email you two keys; you will want to use the "production" key. namedns does not do anything besides manipulate DNS records on domains you already own, so under *Which best describes you?* just pick *I am looking to manage domains via an API*. The process usually takes one or two business days.

### Required global flags

| Key        | Flag               | Description                 |
|:-----------|:-------------------|:----------------------------|
| `username` | `-u`, `--username` | Name.com username           |
| `api-key`  | `-k`, `--api-key`  | Name.com production API key |

#### Optional global flags

| Key       | Flag              | Description                               | Default value           |
|:----------|:------------------|:------------------------------------------|:------------------------|
| `api-url` | `--api-url`       | Name.com base API URL                     | `goname.NameAPIBaseURL` |
|           | `--config`        | Path to a configuration file              | see list below          |
| `output`  | `-o`, `--output`  | Output format: `json`, `table`, or `yaml` | `table`                 |
| `verbose` | `-v`, `--verbose` | Display debugging output                  | `false`                 |

## Configuration

You can save keys to a configuration file so you don't have to specify your username or API key as part of the command. It may be in any format [viper](https://github.com/spf13/viper) recognizes (JSON, TOML, YAML, etc.).

```yaml
username: username
api-key: xxxxx.yyyyy.zzzzz
```

The configuration may be stored in one of these locations, and are searched in the following order:

  - `$XDG_CONFIG_HOME/namedns/.namedns.yaml`.
  - `$XDG_CONFIG_DIRS/namedns/.namedns.yaml`.
  - `$HOME/.namedns.yaml`.
  - `.namedns.yaml`

You may also use environment variables to replace flags; uppercased and prefixed by `NAMEDNS_`.

```bash
NAMEDNS_USERNAME=username NAMEDNS_API_KEY=xxxxx.yyyyy.zzzzz namedns list example.com
```

## Commands

| Command   | Arguments                                  | Description                                  |
|:----------|:-------------------------------------------|:---------------------------------------------|
| `create`  | `<domain> <host> <type> <content> [--ttl]` | Create a DNS record in a domain              |
| `delete`  | `<domain> <record id> [record id ...]`     | Delete one or more DNS records from a domain |
| `list`    | `[domain ...]`                             | List all DNS records for one or more domains |
| `help`    |                                            | Display help and exit                        |
| `version` |                                            | Display version and exit                     |

You may optionally specify a `--ttl` flag (time to live, in seconds) for the `create` and `set` commands. If unspecified, records will be created with a default TTL of 3600 (1 hour). You can set it to as low as 60 according to the API spec, but the production Name.com API does not currently appear to accept a TTL lower than 300.

### Examples

#### General

    # create an A record for www.example.com with the value "192.168.0.1"
    $ namedns create example.com www a "192.168.0.1"
    INFO[0000] record created successfully    content=192.168.0.1 createDate="2017-01-01 12:00:00" name=www.example.com priority=0 recordID=1234 ttl=3600 type=A

    # create a TXT record for www.example.com with the value "hello world!" and a TTL of 300 seconds
    $ namedns create example.com www txt "hello world!" --ttl 300
    INFO[0000] record created successfully    content="hello world!" createDate="2017-01-01 12:00:00" name=www.example.com priority=0 recordID=1235 ttl=300 type=TXT

    # list all DNS records for a domain on an account
    $ namedns list example.com
    # example.com
    1234 example.com A 192.168.0.1 3600
    1235 example.com TXT hello world! 300

    # delete record ID 1234 from example.com
    $ namedns delete example.com 1234
    INFO[0000] record deleted successfully    recordID=1234

#### Usage for [dehydrated] hooks (for Let's Encrypt)

Drop namedns into your `$PATH` and these two functions into `hook.sh`.

```bash
deploy_challenge() {
  local DOMAIN="${1}" TOKEN_FILENAME="${2}" TOKEN_VALUE="${3}"

  SUB_DOMAIN="_acme-challenge.$(sed -E 's/(.*)\.(\w+\.\w+)\.?$/\1/' <<< "$DOMAIN")"
  ROOT_DOMAIN="$(sed -E 's/(.*)\.(\w+\.\w+)\.?$/\2/' <<< "$DOMAIN")"
  [ "$DOMAIN" == "$ROOT_DOMAIN" ] && SUB_DOMAIN='_acme-challenge'

  namedns create "$ROOT_DOMAIN" "$SUB_DOMAIN" txt "$TOKEN_VALUE" --ttl 300 && sleep 90
}

clean_challenge() {
  local DOMAIN="${1}" TOKEN_FILENAME="${2}" TOKEN_VALUE="${3}"

  ROOT_DOMAIN="$(sed -E 's/(.*)\.(\w+\.\w+)\.?$/\2/' <<< "$DOMAIN")"

  record_id=$(namedns list "$ROOT_DOMAIN" | grep "$TOKEN_VALUE" | awk '{print $1}')
  namedns delete "$ROOT_DOMAIN" "$record_id"
}
```

## Building

```bash
dep ensure
go build -ldflags "-X 'github.com/aelindeman/namedns/cmd.Version=$(git describe --tags --candidates=1 --dirty --abbrev=40)'"
./namedns
```

## Author

  - [Alex Lindeman][aelindeman]

[aelindeman]: https://github.com/aelindeman
[dehydrated]: https://github.com/lukas2511/dehydrated
