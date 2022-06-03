package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

type Item struct {
	path string
	size int64
}

type Context struct {
	items    []Item
	capacity int
}

func reverseSort(array []Item) {
	sort.Slice(array, func(a, b int) bool {
		return array[b].size < array[a].size
	})
}

func visitDir(context *Context, dirPath string) {
	items, err := os.ReadDir(dirPath)
	if err != nil {
		return
	}

	for _, item := range items {
		if item.IsDir() {
			visitDir(context, filepath.Join(dirPath, item.Name()))
		} else {
			info, err := item.Info()

			if err == nil {
				item := Item{
					path: filepath.Join(dirPath, item.Name()),
					size: info.Size(),
				}

				if len(context.items) < context.capacity {
					context.items = append(context.items, item)
					reverseSort(context.items)
				} else if item.size > context.items[context.capacity-1].size {
					context.items = append(context.items[:context.capacity-1], item)
					reverseSort(context.items)
				}
			}
		}
	}
}

func main() {
	var initDir = flag.String("dir", ".", "The `path` to start searching in")
	var maxItems = flag.Int("max", 10, "Display `N` largest files")
	flag.Parse()

	context := Context{capacity: *maxItems}
	visitDir(&context, *initDir)

	maxLen := 0
	for _, item := range context.items {
		if len(item.path) > maxLen {
			maxLen = len(item.path)
		}
	}

	for _, item := range context.items {
		fmt.Printf("%-*s => %d\n", maxLen, item.path, item.size)
	}
}
