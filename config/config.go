package config
import "net/url"

var (
	HOST = "https://msghub.eycia.me/"
	BaseUrl *url.URL
)


func init() {
	BaseUrl, _ = url.Parse(HOST)
}