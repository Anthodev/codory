package utils

import (
	"testing"
)

func TestContains(t *testing.T) {
	tests := []struct {
		name   string
		s      string
		substr string
		want   bool
	}{
		{"empty substring", "hello world", "", true},
		{"substring at start", "hello world", "hello", true},
		{"substring at end", "hello world", "world", true},
		{"substring in middle", "hello world", "lo wo", true},
		{"no match", "hello world", "foo", false},
		{"empty string with substring", "", "foo", false},
		{"exact match", "hello", "hello", true},
		{"partial match", "hello", "hell", true},
		{"case sensitive", "Hello", "hello", false},
		{"single character match", "a", "a", true},
		{"single character no match", "a", "b", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Contains(tt.s, tt.substr)
			if got != tt.want {
				t.Errorf("Contains(%q, %q) = %v, want %v", tt.s, tt.substr, got, tt.want)
			}
		})
	}
}

func TestContainsSubstring(t *testing.T) {
	tests := []struct {
		name   string
		s      string
		substr string
		want   bool
	}{
		{"empty substring", "hello world", "", true},
		{"substring at start", "hello world", "hello", true},
		{"substring at end", "hello world", "world", true},
		{"substring in middle", "hello world", "lo wo", true},
		{"no match", "hello world", "foo", false},
		{"empty string with substring", "", "foo", false},
		{"exact match", "hello", "hello", true},
		{"partial match", "hello", "hell", true},
		{"case sensitive", "Hello", "hello", false},
		{"single character match", "a", "a", true},
		{"single character no match", "a", "b", false},
		{"multiple occurrences", "hello hello", "hello", true},
		{"overlapping patterns", "aaaa", "aaa", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ContainsSubstring(tt.s, tt.substr)
			if got != tt.want {
				t.Errorf("ContainsSubstring(%q, %q) = %v, want %v", tt.s, tt.substr, got, tt.want)
			}
		})
	}
}

func TestContainsString(t *testing.T) {
	tests := []struct {
		name   string
		s      string
		substr string
		want   bool
	}{
		{"empty substring", "hello world", "", true},
		{"substring at start", "hello world", "hello", true},
		{"substring at end", "hello world", "world", true},
		{"substring in middle", "hello world", "lo wo", true},
		{"no match", "hello world", "foo", false},
		{"empty string with substring", "", "foo", false},
		{"exact match", "hello", "hello", true},
		{"case sensitive", "Hello", "hello", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ContainsString(tt.s, tt.substr)
			if got != tt.want {
				t.Errorf("ContainsString(%q, %q) = %v, want %v", tt.s, tt.substr, got, tt.want)
			}
		})
	}
}

// BenchmarkContains benchmarks the recursive Contains function
func BenchmarkContains(b *testing.B) {
	s := "hello world, this is a test string with some content"
	substr := "test"
	for i := 0; i < b.N; i++ {
		_ = Contains(s, substr)
	}
}

// BenchmarkContainsSubstring benchmarks the iterative ContainsSubstring function
func BenchmarkContainsSubstring(b *testing.B) {
	s := "hello world, this is a test string with some content"
	substr := "test"
	for i := 0; i < b.N; i++ {
		_ = ContainsSubstring(s, substr)
	}
}

// BenchmarkContainsString benchmarks the standard library strings.Contains
func BenchmarkContainsString(b *testing.B) {
	s := "hello world, this is a test string with some content"
	substr := "test"
	for i := 0; i < b.N; i++ {
		_ = ContainsString(s, substr)
	}
}
