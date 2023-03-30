package azure

import (
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/dns/mgmt/2018-05-01/dns"
	"github.com/Azure/azure-sdk-for-go/services/privatedns/mgmt/2018-09-01/privatedns"
	"github.com/Azure/go-autorest/autorest/to"

	"github.com/stretchr/testify/assert"
)

func Test_parseMxTarget(t *testing.T) {
	type testCase[T interface {
		dns.MxRecord | privatedns.MxRecord
	}] struct {
		name    string
		args    string
		want    T
		wantErr assert.ErrorAssertionFunc
	}

	tests := []testCase[dns.MxRecord]{
		{
			name: "valid mx target",
			args: "10 example.com",
			want: dns.MxRecord{
				Preference: to.Int32Ptr(int32(10)),
				Exchange:   to.StringPtr("example.com"),
			},
			wantErr: assert.NoError,
		},
		{
			name: "valid mx target with a subdomain",
			args: "99 foo-bar.example.com",
			want: dns.MxRecord{
				Preference: to.Int32Ptr(int32(99)),
				Exchange:   to.StringPtr("foo-bar.example.com"),
			},
			wantErr: assert.NoError,
		},
		{
			name:    "invalid mx target with misplaced preference and exchange",
			args:    "example.com 10",
			want:    dns.MxRecord{},
			wantErr: assert.Error,
		},
		{
			name:    "invalid mx target without preference",
			args:    "example.com",
			want:    dns.MxRecord{},
			wantErr: assert.Error,
		},
		{
			name:    "invalid mx target with non numeric preference",
			args:    "aa example.com",
			want:    dns.MxRecord{},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseMxTarget[dns.MxRecord](tt.args)
			if !tt.wantErr(t, err, fmt.Sprintf("parseMxTarget(%v)", tt.args)) {
				return
			}
			assert.Equalf(t, tt.want, got, "parseMxTarget(%v)", tt.args)
		})
	}
}
