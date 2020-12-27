package internal

import (
	"io"
	"os"
	"reflect"
	"testing"
)

func TestParseSiteMap(t *testing.T) {
	singleUrlSitemap, _ := os.Open("../test/data/sitemap_single.xml")
	multiUrlSitemap, _ := os.Open("../test/data/sitemap_multi.xml")
	invalidSitemap, _ := os.Open("../test/data/sitemap_invalid.xml")

	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name:    "Returns single url",
			args:    args{singleUrlSitemap},
			want:    []string{"https://www.google.com"},
			wantErr: false,
		},
		{
			name:    "Returns many urls",
			args:    args{multiUrlSitemap},
			want:    []string{"https://www.google.com", "https://www.yahoo.com", "https://www.bing.com"},
			wantErr: false,
		},
		{
			name:    "Returns err when xml can't be parsed",
			args:    args{invalidSitemap},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseSiteMap(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseSiteMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseSiteMap() got = %v, want %v", got, tt.want)
			}
		})
	}
}

//
//func TestSendRequest(t *testing.T) {
//	type args struct {
//		url    string
//		config *Config
//	}
//	tests := []struct {
//		name    string
//		args    args
//		want    *http.Response
//		wantErr bool
//	}{
//		{
//			name: "Sends an http request",
//			args: args{"localhost:8080", &Config{
//				Auth: NoAuth{},
//			}},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := SendRequest(tt.args.url, tt.args.config)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("SendRequest() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("SendRequest() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
