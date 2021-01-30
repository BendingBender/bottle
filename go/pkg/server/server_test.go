package server

import (
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/remotehack/bottle/pkg/config"
	"github.com/remotehack/bottle/pkg/mocks"
)

func TestServer_getFilename(t *testing.T) {
	r, err := http.NewRequest(http.MethodGet, "https://subdomain.example.com", nil)
	if err != nil {
		t.Fatalf("could not create example request: %s", err)
	}

	type args struct {
		r   *http.Request
		cfg config.Config
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "gets filename if there's a subdomain",
			args: args{
				r: r,
				cfg: config.Config{
					Port: "7889",
					Host: "example.com",
				},
			},
			want:    "subdomain",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPersister := mocks.NewMockPersister(ctrl)
			mockPersister.EXPECT().Write(gomock.Any(), gomock.Any()).AnyTimes()

			s, err := New(tt.args.cfg, mockPersister)
			assert.NoError(t, err)

			got, err := s.getFilename(tt.args.r)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, "", got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
