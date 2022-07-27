Hosted ACME Challenge on Cloudflare for Caddy
===========================

This package contains a DNS provider module for [Caddy](https://github.com/caddyserver/caddy). It can be used to manage hosted ACME challenge records with Cloudflare accounts under a specific zone.

This is useful when you have a Caddy instance/cluster managing domains across different DNS providers, and you need a single point of entry for ACME challenges.

## How to use

If you want to host acme challenges for `secure.example.com` (or its wildcard variant), add a CNAME record `_acme-challenge.secure.example.com` pointing to `secure.example.com.acmehostedexample.com`, where `acmehostedexample.com` is the zone in the config and in Cloudflare.

## Caddy module name

```
dns.providers.challenger
```

## Config examples

To use this module for the ACME DNS challenge, [configure the ACME issuer in your Caddy JSON](https://caddyserver.com/docs/json/apps/tls/automation/policies/issuer/acme/) like so:

```json
{
	"module": "acme",
	"challenges": {
		"dns": {
			"provider": {
				"name": "challenger",
				"api_token": "{env.CF_API_TOKEN}",
                "root_zone": "acmehostedexample.com"
			}
		}
	}
}
```

or with the Caddyfile:

```
tls {
	dns challenger {env.CF_API_TOKEN} acmehostedexample.com
}
```

You can replace `{env.CF_API_TOKEN}` with the actual auth token if you prefer to put it directly in your config instead of an environment variable.

`acmehostedexample.com` is the zone responsible for hosting acme challenges for other zones.

## Authenticating

See [the associated README in the libdns package](https://github.com/libdns/cloudflare) for important information about credentials.

**NOTE**: If migrating from Caddy v1, you will need to change from using a Cloudflare API Key to a scoped API Token. Please see link above for more information.

## Troubleshooting

### Error: `Invalid request headers`

If providing your API token via an ENV var which is accidentally not set/available when running Caddy, you'll receive this error from Cloudflare.

Double check that Caddy has access to a valid CF API token.

### Error: `timed out waiting for record to fully propagate`

Some environments may have trouble querying the `_acme-challenge` TXT record from Cloudflare. Verify in the Cloudflare dashboard that the temporary record is being created.

If the record does exist, your DNS resolver may be caching an earlier response before the record was valid. You can instead configure Caddy to use an alternative DNS resolver such as [Cloudflare's official `1.1.1.1`](https://www.cloudflare.com/en-gb/learning/dns/what-is-1.1.1.1/).

Add a custom `resolver` to the [`tls` directive](https://caddyserver.com/docs/caddyfile/directives/tls):

```
tls {
  dns challenger {env.CF_API_TOKEN} acmehostedexample.com
  resolvers 1.1.1.1
}
```

Or with Caddy JSON to the `acme` module: [`challenges.dns.provider.resolvers: ["1.1.1.1"]`](https://caddyserver.com/docs/json/apps/tls/automation/policies/issuer/acme/challenges/dns/resolvers/).

## Credits

This package is based on [libdns/cloudflare](https://github.com/libdns/cloudflare) and [caddy-dns/cloudflare](https://github.com/caddy-dns/cloudflare)