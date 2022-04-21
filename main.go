package main

import (
	"fmt"
	"os"
	"sync"
	"strings"
	fp "path/filepath"
)

type Counter struct {
	Value int
	Mutex sync.Mutex
}

// COUNT
func count(file string, c *Counter, wg *sync.WaitGroup) {
	defer wg.Done()

	contents, err := os.ReadFile(file)
	if err != nil {
		fmt.Println(file, err)
		return
	}

	lines := strings.Split(string(contents), "\n")

	c.Mutex.Lock()
	c.Value += len(lines)
	c.Mutex.Unlock()
}

// READ
func read(dir string, c *Counter, wg *sync.WaitGroup) {
	defer wg.Done()

	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println(dir, err)
		return
	}

	for _, file := range files {
		path := fp.Join(dir, file.Name())
		
		if file.IsDir() {
			wg.Add(1)
			go read(path, c, wg)
		} else if fp.Ext(file.Name()) == ".go" {
			wg.Add(1)
			go count(path, c, wg)
		}
	}

}

// MAIN
func main() {
	var c Counter
	var wg sync.WaitGroup
	
	wg.Add(1)
	go read(".", &c, &wg)

	wg.Wait()
	fmt.Println(c.Value)
}