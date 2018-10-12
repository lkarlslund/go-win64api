// +build windows,amd64

package winapi

import (
	"testing"
)

func TestUpdatesPending(t *testing.T) {
	got, err := UpdatesPending()
	if err != nil {
		t.Errorf("UpdatesPending() error = %v, wantErr %v", err, tt.wantErr)
		return
	}
	if got == nil {
		t.Errorf("UpdatesPending() returned nil object")
		return
	}
}
