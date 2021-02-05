package kustomizev3

import (
	"fmt"
	"reflect"
	"testing"

	iacloaderror "github.com/accurics/terrascan/pkg/iac-providers/iac-load-error"
	"github.com/accurics/terrascan/pkg/iac-providers/output"
	"github.com/accurics/terrascan/pkg/utils"
	"go.uber.org/zap"
)

var errorReadKustomize = fmt.Errorf("error while reading the file %s %v", "kustomization.yaml", zap.Error(utils.ErrYamlFileEmpty))

func TestLoadIacDir(t *testing.T) {

	table := []struct {
		name          string
		dirPath       string
		kustomize     KustomizeV3
		want          output.AllResourceConfigs
		wantErr       error
		resourceCount int
	}{
		{
			name:          "invalid dirPath",
			dirPath:       "not-there",
			kustomize:     KustomizeV3{},
			wantErr:       errorDirRead,
			resourceCount: 0,
		},
		{
			name:          "simple-deployment",
			dirPath:       "./testdata/simple-deployment",
			kustomize:     KustomizeV3{},
			wantErr:       nil,
			resourceCount: 4,
		},
		{
			name:          "multibases",
			dirPath:       "./testdata/multibases/base",
			kustomize:     KustomizeV3{},
			wantErr:       nil,
			resourceCount: 2,
		},
		{
			name:          "multibases",
			dirPath:       "./testdata/multibases/dev",
			kustomize:     KustomizeV3{},
			wantErr:       nil,
			resourceCount: 2,
		},
		{
			name:          "multibases",
			dirPath:       "./testdata/multibases/prod",
			kustomize:     KustomizeV3{},
			wantErr:       nil,
			resourceCount: 2,
		},

		{
			name:          "multibases",
			dirPath:       "./testdata/multibases/stage",
			kustomize:     KustomizeV3{},
			wantErr:       nil,
			resourceCount: 2,
		},
		{
			name:          "multibases",
			dirPath:       "./testdata/multibases",
			kustomize:     KustomizeV3{},
			wantErr:       nil,
			resourceCount: 4,
		},
		{
			name:          "no-kustomize-directory",
			dirPath:       "./testdata/no-kustomizefile",
			kustomize:     KustomizeV3{},
			wantErr:       errorKustomizeNotFound,
			resourceCount: 0,
		},
		{
			name:          "kustomize-file-empty",
			dirPath:       "./testdata/kustomize-file-empty",
			kustomize:     KustomizeV3{},
			wantErr:       utils.ErrYamlFileEmpty,
			resourceCount: 0,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			resourceMap, gotErr := tt.kustomize.LoadIacDir(tt.dirPath)
			if gotErr != nil {
				if e, ok := gotErr.(*iacloaderror.LoadError); !ok || e.Err != tt.wantErr {
					t.Errorf("TestLoadIacDir()= gotErr: '%v', wantErr: '%v'", gotErr, tt.wantErr)
				}
			} else {
				if gotErr != tt.wantErr {
					t.Errorf("TestLoadIacDir()= gotErr: '%v', wantErr: '%v'", gotErr, tt.wantErr)
				}
			}

			resCount := utils.GetResourceCount(resourceMap)
			if resCount != tt.resourceCount {
				t.Errorf("TestLoadIacDir()= resource count (%d) does not match expected (%d)", resCount, tt.resourceCount)
			}
		})
	}

}

func TestLoadKustomize(t *testing.T) {
	kustomizeYaml := "kustomization.yaml"
	kustomizeYml := "kustomization.yml"

	table := []struct {
		name     string
		basepath string
		filename string
		want     output.AllResourceConfigs
		wantErr  error
	}{
		{
			name:     "simple-deployment",
			basepath: "./testdata/simple-deployment",
			filename: kustomizeYaml,
			wantErr:  nil,
		},
		{
			name:     "multibases",
			basepath: "./testdata/multibases",
			filename: kustomizeYaml,
			wantErr:  nil,
		},
		{
			name:     "multibases/base",
			basepath: "./testdata/multibases/base",
			filename: kustomizeYml,
			wantErr:  nil,
		},
		{
			name:     "multibases/dev",
			basepath: "./testdata/multibases/dev",
			filename: kustomizeYaml,
			wantErr:  nil,
		},
		{
			name:     "multibases/prod",
			basepath: "./testdata/multibases/prod",
			filename: kustomizeYaml,
			wantErr:  nil,
		},
		{
			name:     "multibases/stage",
			basepath: "./testdata/multibases/stage",
			filename: kustomizeYaml,
			wantErr:  nil,
		},
		{
			name:     "multibases/zero-violation-base",
			basepath: "./testdata/multibases/zero-violation-base",
			filename: kustomizeYaml,
			wantErr:  nil,
		},
		{
			name:     "erroneous-pod",
			basepath: "./testdata/erroneous-pod",
			filename: kustomizeYaml,
			wantErr:  errorFromKustomize,
		},
		{
			name:     "erroneous-deployment",
			basepath: "./testdata/erroneous-deployment/",
			filename: kustomizeYaml,
			wantErr:  errorFromKustomize,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			_, gotErr := LoadKustomize(tt.basepath, tt.filename)
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("unexpected error; gotErr: '%v', wantErr: '%v'", gotErr, tt.wantErr)
			}
		})
	}
}
