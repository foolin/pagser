package pagser

import "testing"

func TestToInt32SliceE(t *testing.T) {
	_, err := toFloat32SliceE(nil)
	if err == nil {
		t.Fatal("nil not return error")
	}

	_, err = toFloat32SliceE(1)
	if err == nil {
		t.Fatal("1 not return error")
	}

	arrs := []int32{1, 2, 3}
	out, err := toFloat32SliceE(arrs)
	t.Logf("out: %v, error: %v", out, err)

}
