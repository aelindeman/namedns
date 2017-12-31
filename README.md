# namedns

A command-line utility to manipulate Name.com DNS records

## Install

Grab the latest binary from the [releases](https://github.com/aelindeman/namedns/releases) page and drop it in your `$PATH`:

    curl -o namedns "https://github.com/aelindeman/namedns/releases/download/v1.0/namedns-$(uname -s)-$(uname -m)"
    chmod +x namedns
    mv namedns /usr/local/bin/namedns

## Configuration

Configuration can be passed with flags on the command line, or can be stored in `$HOME/.namedns.yaml`. A username and API key are required.

| Key        | Flag               | Description                 |
|:-----------|:-------------------|:----------------------------|
| `username` | `-u`, `--username` | Name.com username           |
| `api-key`  | `-k`, `--api-key`  | Name.com production API key |

You can get an API key by going to <https://name.com/reseller/apply> and filling out the form. namedns does not do anything besides manipulate DNS records on domains you already own, so under *Which best describes you?* just pick *I am looking to manage domains via an API*. The process usually takes one or two business days.

## Usage

| Command   | Arguments                                  | Description                                  |
|:----------|:-------------------------------------------|:---------------------------------------------|
| `create`  | `<domain> <host> <type> <content> [flags]` | Create a DNS record in a domain              |
| `delete`  | `<domain> <record id> [record id ...]`     | Delete one or more DNS records from a domain |
| `list`    | `[domain ...]`                             | List all DNS records for one or more domains |
| `help`    |                                            | Display help and exit                        |
| `version` |                                            | Display version and exit                     |

The `create` command can take an optional `--ttl` flag (time to live, in seconds). If unspecified, records will be created with a TTL of 3600 (1 hour). You can set it to as low as 60 according to the API spec, but the Name.com API does not currently appear to accept any TTL lower than 300.

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

#### Use in [dehydrated] hook (for Let's Encrypt)

Drop namedns into your `$PATH` and these two functions into `hook.sh`.

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

## Building

    glide install
    go build
    ./namedns

## Author

  - [Alex Lindeman][aelindeman]

[aelindeman]: https://github.com/aelindeman
[dehydrated]: https://github.com/lukas2511/dehydrated
