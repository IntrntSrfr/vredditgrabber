package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
)

func main() {

	for {
		inp := ""
		fmt.Scan(&inp)
		err := getRedditVideo(inp+".json", "./gamer.mp4")
		if err != nil {
			fmt.Println(err)
		}
	}

}

func getRedditVideo(url, filepath string) error {

	p, err := getPost(url)
	if err != nil {
		return err
	}

	if !p.IsVideo || p.Domain != "v.redd.it" {
		return errors.New("this isnt gonna work")
	}

	video := "./video.mp4"
	/*
		if p.Media.RedditVideo.IsGif {
			video = "./video.gif"
		}
	*/
	audio := "./audio.mp4"

	err = getVideo(video, p.Media.RedditVideo.FallbackURL)
	if err != nil {
		return err
	}

	if !p.Media.RedditVideo.IsGif {

		err = getAudio(audio, p.URL+"/audio")
		if err != nil {
			return err
		}

		err = combine(video, audio, filepath)
		if err != nil {
			fmt.Println(err)
		}
		os.Remove(audio)
	}

	os.Remove(video)

	return nil
}

func getPost(url string) (*Post, error) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", "weed")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	p := []PostDetailed{}
	err = json.NewDecoder(res.Body).Decode(&p)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return p[0].Data.Children[0].Post, nil
}

func getVideo(name, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(name)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func getAudio(name, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(name)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func combine(video, audio, filepath string) error {
	return exec.Command("ffmpeg", "-i", video, "-i", audio, "-c", "copy", filepath).Run()
}

type Post struct {
	Title                 string `json:"title"`
	SubredditNamePrefixed string `json:"subreddit_name_prefixed"`
	Domain                string `json:"domain"`
	Author                string `json:"author"`

	CrosspostParentList []struct {
		*Post
	} `json:"crosspost_parent_list"`
	Media struct {
		RedditVideo struct {
			FallbackURL string `json:"fallback_url"`
			DashURL     string `json:"dash_url"`
			IsGif       bool   `json:"is_gif"`
		} `json:"reddit_video"`
	} `json:"media"`
	URL     string `json:"url"`
	IsVideo bool   `json:"is_video"`
}

type PostDetailed struct {
	Kind string `json:"kind"`
	Data struct {
		Modhash  string `json:"modhash"`
		Dist     int    `json:"dist"`
		Children []struct {
			*Post `json:"data"`
		} `json:"children"`
	} `json:"data"`
}
