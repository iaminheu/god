package stringx

import (
	"fmt"
	"testing"
)

func TestBytesToSize(t *testing.T) {
	fmt.Println("435 Bytes", BytesToSize(435))
	fmt.Println("3.3 KB", BytesToSize(3398))
	fmt.Println("479 KB", BytesToSize(490398))
	fmt.Println("6.2 MB", BytesToSize(6544528))
	fmt.Println("22 MB", BytesToSize(23483023))
	fmt.Println("3.7 GB", BytesToSize(3984578493))
	fmt.Println("28 GB", BytesToSize(30498505889))
	fmt.Println("8.4 PB", BytesToSize(9485039485039445))
}

func TestContains(t *testing.T) {
	fmt.Println("435 Bytes", Bytes2Size(435))
	fmt.Println("3.3 KB", Bytes2Size(3398))
	fmt.Println("479 KB", Bytes2Size(490398))
	fmt.Println("6.2 MB", Bytes2Size(6544528))
	fmt.Println("22 MB", Bytes2Size(23483023))
	fmt.Println("3.7 GB", Bytes2Size(3984578493))
	fmt.Println("28 GB", Bytes2Size(30498505889))
	fmt.Println("8.4 PB", Bytes2Size(9485039485039445))
}
