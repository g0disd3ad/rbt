package dictionary

import (
	"bufio"
	"cmp"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/g0disd3ad/rbt/internal/tree"
)

var (
	ErrorInvalidFormat = errors.New("Invalid format")
	ErrorWordNotFound  = errors.New("Word is not found")
)

type TreeStorage[K cmp.Ordered, V any] interface {
	InsertOrGet(key K) *tree.Node[K, V]
	Remove(key K) error
	Search(key K) (*tree.Node[K, V], bool)
	InOrderWalk(func(key K, value V) bool)
	Height() int
	IsValid() bool
}

type RBTStorage struct {
	t TreeStorage[string, []string]
}

func NewRBTStorage() *RBTStorage {
	return &RBTStorage{
		t: tree.NewTree[string, []string](),
	}
}

func (s *RBTStorage) Insert(key, translation string) error {
	node := s.t.InsertOrGet(key)
	for _, v := range node.Value {
		if v == translation {
			return nil
		}
	}
	node.Value = append(node.Value, translation)
	return nil
}

func (s *RBTStorage) Remove(key string) error {
	return s.t.Remove(key)
}

func (s *RBTStorage) Search(key string) ([]string, error) {
	node, found := s.t.Search(key)
	if !found || len(node.Value) == 0 {
		return nil, ErrorWordNotFound
	}
	return node.Value, nil
}

func (s *RBTStorage) InOrderWalk(fn func(key string, translations []string) bool) {
	s.t.InOrderWalk(fn)
}

func (s *RBTStorage) Height() int {
	return s.t.Height()
}

func (s *RBTStorage) IsValid() bool {
	return s.t.IsValid()
}

type Storage interface {
	Insert(key, translation string) error
	Remove(key string) error
	Search(key string) ([]string, error)
	InOrderWalk(func(key string, translations []string) bool)
	Height() int
	IsValid() bool
}

type Dictionary struct {
	data Storage
}

func NewDictionary(s Storage) *Dictionary {
	return &Dictionary{
		data: s,
	}
}

func (d *Dictionary) Insert(key, translation string) error {
	key = strings.TrimSpace(strings.ToLower(key))
	translation = strings.TrimSpace(strings.ToLower(translation))
	return d.data.Insert(key, translation)
}

func (d *Dictionary) Remove(key string) error {
	key = strings.TrimSpace(strings.ToLower(key))
	return d.data.Remove(key)
}

func (d *Dictionary) Search(key string) ([]string, error) {
	key = strings.TrimSpace(strings.ToLower(key))
	return d.data.Search(key)
}

func (d *Dictionary) InOrderWalk(fn func(key string, translations []string) bool) {
	d.data.InOrderWalk(fn)
}

func (d *Dictionary) LoadFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("Could not open the file %s", filename)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
		line := scanner.Text()
		if line == "" {
			continue
		}
		parts := strings.Split(line, "-")
		if len(parts) == 2 {
			eng := strings.TrimSpace(parts[0])
			rus := strings.TrimSpace(parts[1])
			if eng != "" && rus != "" {
				if err := d.Insert(eng, rus); err != nil {
					fmt.Fprintf(os.Stderr, "--Warning: line %d: %v\n", lineCount, err)
				}
			} else {
				fmt.Fprintf(os.Stderr, "--Warning: line %d is skipped - an empty word\n", lineCount)
			}
		} else {
			fmt.Fprintf(os.Stderr, "--Warning: line %d is skipped - invalid format\n", lineCount)
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("Error while reading the file - %v", err)
	}
	fmt.Printf("Dictionary is loaded from file %s successfully\n", filename)
	return nil
}

func (d *Dictionary) SaveToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("Error saving to the file %s: %v", filename, err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	d.data.InOrderWalk(func(key string, translations []string) bool {
		for _, tr := range translations {
			_, err = fmt.Fprintf(writer, "%s - %s\n", key, tr)
			if err != nil {
				return false
			}
		}
		return true
	})
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("Error saving buffer to the file %s: %v", filename, err)
	}
	fmt.Printf("Dictionary is saved to %s\n", filename)
	return nil
}

func (d *Dictionary) Print() {
	fmt.Println("\n~--------------------------------------~")
	lineCount := 1
	d.data.InOrderWalk(func(key string, translations []string) bool {
		fmt.Printf("%d. %s : %s\n", lineCount, key, strings.Join(translations, ", "))
		lineCount++
		return true
	})
	fmt.Println("~--------------------------------------~")
}

func (d *Dictionary) GetHeight() int {
	return d.data.Height()
}

func (d *Dictionary) IsValidTree() bool {
	return d.data.IsValid()
}
