package geodata

import v2router "github.com/v2fly/v2ray-core/v5/app/router/routercommon"

type GeodataLoader interface {
	LoadIP(filename, country string) ([]*v2router.CIDR, error)
	LoadSite(filename, list string) ([]*v2router.Domain, error)
	LoadGeoIP(country string) ([]*v2router.CIDR, error)
	LoadGeoSite(list string) ([]*v2router.Domain, error)
}
