package pagser

import "testing"

func TestToInt32SliceE(t *testing.T) {
	_, err := toInt32SliceE(nil)
	if err == nil {
		t.Fatal("nil not return error")
	}

	_, err = toInt32SliceE(1)
	if err == nil {
		t.Fatal("1 not return error")
	}

	list := []int32{1, 2, 3}
	out, err := toInt32SliceE(list)
	t.Logf("out: %v, error: %v", out, err)
}

func TestToInt64SliceE(t *testing.T) {
	_, err := toInt64SliceE(nil)
	if err == nil {
		t.Fatal("nil not return error")
	}

	_, err = toInt64SliceE(1)
	if err == nil {
		t.Fatal("1 not return error")
	}

	list := []int64{1, 2, 3}
	out, err := toInt64SliceE(list)
	t.Logf("out: %v, error: %v", out, err)
}

func TestToFloat32SliceE(t *testing.T) {
	_, err := toFloat32SliceE(nil)
	if err == nil {
		t.Fatal("nil not return error")
	}

	_, err = toFloat32SliceE(1)
	if err == nil {
		t.Fatal("1 not return error")
	}

	list := []float32{1, 2, 3}
	out, err := toFloat32SliceE(list)
	t.Logf("out: %v, error: %v", out, err)
}

func TestToFloat64SliceE(t *testing.T) {
	_, err := toFloat64SliceE(nil)
	if err == nil {
		t.Fatal("nil not return error")
	}

	_, err = toFloat64SliceE(1)
	if err == nil {
		t.Fatal("1 not return error")
	}

	list := []float64{1, 2, 3}
	out, err := toFloat64SliceE(list)
	t.Logf("out: %v, error: %v", out, err)
}
