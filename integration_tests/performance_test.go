package integration_tests

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func BenchmarkStreamWrite(t *testing.B) {
	buf := make([]byte, 1024*1024)
	t.SetBytes(int64(len(buf)))

	file, err := os.Create(defaultPlainDir + "BenchmarkWrite")
	if err != nil {
		t.FailNow()
	}

	t.ResetTimer()
	var i int
	for i = 0; i < t.N; i++ {
		written, err := file.Write(buf)
		if err != nil {
			fmt.Printf("err=\"%s\", written=%d\n", err.Error(), written)
			t.FailNow()
		}
	}
	file.Close()
}

func BenchmarkStreamRead(t *testing.B) {
	buf := make([]byte, 1024*1024)
	t.SetBytes(int64(len(buf)))

	fn := defaultPlainDir + "BenchmarkWrite"
	fi, err := os.Stat(fn)
	if err != nil {
		t.Fatal(err)
	}
	mb := int(fi.Size() / 1024 / 1024)

	if t.N > mb {
		// Grow file so we can satisfy the test
		//fmt.Printf("Growing file to %d MB... ", t.N)
		f2, err := os.OpenFile(fn, os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println(err)
			t.FailNow()
		}
		for h := 0; h < t.N-mb; h++ {
			_, err = f2.Write(buf)
			if err != nil {
				fmt.Println(err)
				t.FailNow()
			}
		}
		f2.Close()
		//fmt.Printf("done\n")
	}

	file, err := os.Open(fn)
	if err != nil {
		t.FailNow()
	}
	t.ResetTimer()
	var i int
	for i = 0; i < t.N; i++ {
		_, err := file.Read(buf)
		if err == io.EOF {
			fmt.Printf("Test file too small\n")
			t.SkipNow()
		} else if err != nil {
			fmt.Println(err)
			t.FailNow()
		}
	}
	file.Close()
}

// createFiles - create "count" files of size "size" bytes each
func createFiles(t *testing.B, count int, size int) {
	dir := fmt.Sprintf("%s/createFiles_%d_%d", defaultPlainDir, count, size)
	err := os.Mkdir(dir, 0777)
	if err != nil {
		t.Fatal(err)
	}
	buf := make([]byte, size)
	t.SetBytes(int64(len(buf)))
	t.ResetTimer()
	var i int
	for i = 0; i < count; i++ {
		file := fmt.Sprintf("%s/%d", dir, i)
		if size > 0 {
			err = ioutil.WriteFile(file, buf, 0666)
		} else {
			var fh *os.File
			fh, err = os.Create(file)
			fh.Close()
		}
		if err != nil {
			t.Fatal(err)
		}
	}
	t.StopTimer()
	os.RemoveAll(dir)
}

func BenchmarkCreate0B(t *testing.B) {
	createFiles(t, t.N, 0)
}

func BenchmarkCreate1B(t *testing.B) {
	createFiles(t, t.N, 1)
}

func BenchmarkCreate100B(t *testing.B) {
	createFiles(t, t.N, 100)
}

func BenchmarkCreate4kB(t *testing.B) {
	createFiles(t, t.N, 4*1024)
}

func BenchmarkCreate10kB(t *testing.B) {
	createFiles(t, t.N, 10*1024)
}
