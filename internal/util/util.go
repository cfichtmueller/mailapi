package util

func Nvl(in, fallback string) string {
	if in == "" {
		return fallback
	}
	return in
}
