package code

type Resolver interface {
	SupportedCodes() (localeNames map[string]string)
	ResolveCode(c string) (localeName string, found bool)
	LookupName(n string) (code string, found bool)
}
