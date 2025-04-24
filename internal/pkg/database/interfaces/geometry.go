package interfaces

type GeometryIF interface {
	GetGeoFromText(wkt string) (string, error)
	GetWktFromGeo(geo string) (string, error)
	GetGeoJson(geo string) (string, error)
	GetBufferFromGeoHex(geo string) (string, error)
	GeometryEncode(geometry string) string
	GetWKTFromWKB(wkb string) (string, error)
}
