package main

import (
	"bufio"
	"os"
	"testing"
	"time"

	"encoding/json"

	"github.com/mattermost/mattermost-server/model"
)

func loadPosts(t *testing.T) []model.Post {
	jsonfile, err := os.Open("test-fixtures/posts.json")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	scanner := bufio.NewScanner(jsonfile)

	posts := []model.Post{}

	for scanner.Scan() {
		post := model.Post{}
		err := json.Unmarshal([]byte(scanner.Text()), &post)
		if err != nil {
			t.Fatalf("err: %s", err)
		}
		posts = append(posts, post)

	}

	return posts
}

func TestMatchMessage(t *testing.T) {
	checks := []struct {
		message  string
		expected bool
	}{
		{"123", false},
		{"test test", false},
		{"yaml yamluaieuaidtn", true},
		{"faooudtainedutraine", false},
		{"dies ist ein yaml", true},
		{"fooyamlbar", false},
		{"yaml", true},
		{".yaml", false},
		{"Ya Ml", false},
		{"yAmL", true},
	}

	for _, check := range checks {
		actual := matchMessage(check.message)
		if actual != check.expected {
			t.Fatalf("%s: %t != %t", check.message, actual, check.expected)
		}
	}
}

func TestMatchFrequency(t *testing.T) {
	checks := []struct {
		start, end time.Time
		expected   bool
	}{
		{time.Unix(42, 0), time.Unix(43, 0), true},
		{time.Unix(42, 0), time.Unix(50, 0), false},
	}

	for _, check := range checks {
		actual := matchFrequency(check.start, check.end)
		if actual != check.expected {
			t.Fatalf("%s: %t != %t", check.end.Sub(check.start), actual, check.expected)
		}
	}
}
