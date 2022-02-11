package components

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGet(t *testing.T) {
	type args struct {
		filename     string
		unstructured bool
		data         interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Set worker replicas",
			args: args{
				filename:     "worker-deployment.yaml",
				unstructured: false,
				data:         map[string]string{"workerReplicas": "5"},
			},
			wantErr: false,
		},
		{
			name: "Set starter starterRate",
			args: args{
				filename:     "simple-starter-deployment.yaml",
				unstructured: false,
				data:         map[string]string{"simpleStarterRate": "500"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		err := LoadTemplates()
		assert.NoError(t, err)
		t.Run(tt.name, func(t *testing.T) {
			obj, err := Get(tt.args.filename, tt.args.unstructured, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.NotNil(t, obj)
		})
	}
}
