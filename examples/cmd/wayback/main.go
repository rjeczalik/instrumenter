package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/rjeczalik/instrumenter"
)

func curl(url string, v interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return instrumenter.Errorf("error getting %s: %s", url, err)
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(v)
	if err != nil {
		return instrumenter.Error("invalid json payload: " + err.Error())
	}
	return nil
}

type waybackReq struct {
	ArchivedSnapshots struct {
		Closest struct {
			URL string `json:"url"`
		} `json:"closest"`
	} `json:"archived_snapshots"`
}

func wayback(url string, t time.Time) (string, error) {
	const layout = "20060102150405"
	const api = "http://archive.org/wayback/available?url=%s&timestamp=%s"

	var resp waybackReq
	err := curl(fmt.Sprintf(api, url, t.Format(layout)), &resp)
	if err != nil {
		return "", instrumenter.Errorf("error calling wayback api: %s", err)
	}
	snapshot := resp.ArchivedSnapshots.Closest.URL
	if snapshot == "" {
		return "", instrumenter.Error("no snapshot found")
	}
	return snapshot, nil
}

var back = flag.Int("b", 0, "How old snapshot should be (in days).")

func main() {
	flag.Parse()

	timestamp := time.Now().UTC()
	if *back != 0 {
		timestamp = timestamp.AddDate(0, 0, -*back)
	}

	snapshot, err := wayback(flag.Arg(0), timestamp)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(snapshot)
}
