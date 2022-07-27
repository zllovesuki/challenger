package challenger

import (
	"kon.nect.sh/challenger/cloudflare"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
)

// Provider wraps the provider implementation as a Caddy module.
type Provider struct{ *cloudflare.Provider }

func init() {
	caddy.RegisterModule(Provider{})
}

// CaddyModule returns the Caddy module information.
func (Provider) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "dns.providers.challenger",
		New: func() caddy.Module { return &Provider{new(cloudflare.Provider)} },
	}
}

// Before using the provider config, resolve placeholders in the API token and RootZone
// Implements caddy.Provisioner.
func (p *Provider) Provision(ctx caddy.Context) error {
	r := caddy.NewReplacer()
	p.Provider.APIToken = r.ReplaceAll(p.Provider.APIToken, "")
	p.Provider.RootZone = r.ReplaceAll(p.Provider.RootZone, "")
	return nil
}

// UnmarshalCaddyfile sets up the DNS provider from Caddyfile tokens. Syntax:
//
// challenger [<api_token>] [zone_hosting_challenges] {
//     api_token <api_token>
//     root_zone <zone_hosting_challenges>
// }
//
// Expansion of placeholders in the API token is left to the JSON config caddy.Provisioner (above).
func (p *Provider) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		if d.NextArg() {
			p.Provider.APIToken = d.Val()
		}
		if d.NextArg() {
			p.Provider.RootZone = d.Val()
		}
		if d.NextArg() {
			return d.ArgErr()
		}
		for nesting := d.Nesting(); d.NextBlock(nesting); {
			switch d.Val() {
			case "api_token":
				if p.Provider.APIToken != "" {
					return d.Err("API token already set")
				}
				p.Provider.APIToken = d.Val()
			case "root_zone":
				if p.Provider.RootZone != "" {
					return d.Err("Root Zone already set")
				}
				p.Provider.RootZone = d.Val()
				if d.NextArg() {
					return d.ArgErr()
				}
			default:
				return d.Errf("unrecognized subdirective '%s'", d.Val())
			}
		}
	}
	if p.Provider.APIToken == "" {
		return d.Err("missing API token")
	}
	if p.Provider.RootZone == "" {
		return d.Err("missing root zone")
	}
	return nil
}

// Interface guards
var (
	_ caddyfile.Unmarshaler = (*Provider)(nil)
	_ caddy.Provisioner     = (*Provider)(nil)
	_ caddy.Validator       = (*Provider)(nil)
)
