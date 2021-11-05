package proxy

type Definition struct {
	ListenPath  string `bson:"listen_path" json:"listen_path" mapstructure:"listen_path" valid:"required"`
	UpstreamURL string `bson:"upstream_url" json:"upstream_url" mapstructure:"upstream_url" valid:"url,required"`
	Methods     []string `bson:"methods" json:"methods" mapstructure:"methods" valid:"required"`
}
