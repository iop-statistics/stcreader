package stc

import (
	"os"
	"testing"
)

const (
	TestFile = "1.stc"
)


func TestNewReader(t *testing.T) {
	f, err := os.Open(TestFile)
	if err != nil {
		t.Fatal(err)
	}

	rd, err := NewReader(f)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Read file with code %d, data count: %d, type count: %d\n", rd.Header.Code, rd.Header.DataCount, rd.Header.TypeCount)
}

func TestReader_ReadRaw(t *testing.T) {
	f, err := os.Open(TestFile)
	if err != nil {
		t.Fatal(err)
	}

	rd, err := NewReader(f)
	if err != nil {
		t.Fatal(err)
	}


	r, err := rd.ReadRaw()
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range r.Data {
		t.Logf("%+v\n", v)
	}
}

func TestReader_ReadAllRaw(t *testing.T) {
	f, err := os.Open(TestFile)
	if err != nil {
		t.Fatal(err)
	}

	rd, err := NewReader(f)
	if err != nil {
		t.Fatal(err)
	}


	r, err := rd.ReadAllRaw()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Get %d row(s).\n", len(r))
}

func TestReader_Read(t *testing.T) {
	var st struct {
		I1 int32
		I2 int32
		S1 string
		S2 string
		S3 string
		S4 string
		I3 int32
		I4 int32
		S5 string
		S6 string
		I5 int32
	}

	f, err := os.Open(TestFile)
	if err != nil {
		t.Fatal(err)
	}

	rd, err := NewReader(f)
	if err != nil {
		t.Fatal(err)
	}

	err = rd.Read(&st)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v\n", st)
}

func TestReader_ReadAll(t *testing.T) {
	var st []struct {
		I1 int32
		I2 int32
		S1 string
		S2 string
		S3 string
		S4 string
		I3 int32
		I4 int32
		S5 string
		S6 string
		I5 int32
	}

	f, err := os.Open(TestFile)
	if err != nil {
		t.Fatal(err)
	}

	rd, err := NewReader(f)
	if err != nil {
		t.Fatal(err)
	}

	err = rd.ReadAll(&st)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Read %d data(s).", len(st))
}
