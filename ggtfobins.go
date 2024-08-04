package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gookit/color"
)

var (
	badass  = color.HEX("#bada55")
	hotpink = color.HEX("#f06")
)

type Flags struct {
	bins    string
	exploit string
}

func main() {
	flag.Usage = func() {
		printBanner()
	}

	flags := getFlags()
	flag.Parse()

	bins := flags.bins
	exploit := flags.exploit
	err := validateRequiredFlagValues(bins, exploit)
	if err != nil {
		printBanner()
		color.Red.Println(err)
		return
	}

	printFlagsBanner(exploit, bins)

	binsList := strings.Split(bins, ",")
	for _, bin := range binsList {
		trimmedBin := strings.TrimSpace(bin)
		url := fmt.Sprintf("https://gtfobins.github.io/gtfobins/%s", trimmedBin)
		err = printBins(url, trimmedBin, exploit)
		if err != nil {
			log.Println(err)
		}
	}

	printCredits()
}

func printBanner() {
	badass.Print("\n ______     ______     ______   ______   ______     ______     __     __   __     ______    \n/\\  ___\\   /\\  ___\\   /\\__  _\\ /\\  ___\\ /\\  __ \\   /\\  == \\   /\\ \\   /\\ \"-.\\ \\   /\\  ___\\   \n\\ \\ \\__ \\  \\ \\ \\__ \\  \\/_/\\ \\/ \\ \\  __\\ \\ \\ \\/\\ \\  \\ \\  __<   \\ \\ \\  \\ \\ \\-.  \\  \\ \\___  \\  \n \\ \\_____\\  \\ \\_____\\    \\ \\_\\  \\ \\_\\    \\ \\_____\\  \\ \\_____\\  \\ \\_\\  \\ \\_\\\\\"\\_\\  \\/\\_____\\ \n  \\/_____/   \\/_____/     \\/_/   \\/_/     \\/_____/   \\/_____/   \\/_/   \\/_/ \\/_/   \\/_____/")
	fmt.Print("\n\n")
	color.White.Print("Get info about a given exploit for given bins\n\n")
	badass.Print("Usage: ggtfobins  --exploit suid --bins bash,cat\n\n")
	color.White.Print("Available exploits:\n\n- bind-shell\n- capabilities\n- bin\n- file-download\n- file-read\n- file-upload\n- file-write\n- library-load\n- limited-suid\n- non-interactive-bind-shell\n- non-interactive-reverse-shell\n- reverse-shell\n- shell\n- sudo\n- suid\n\n")
}

func getFlags() Flags {
	binsPtr := flag.String("bins", "", "Comma-separated list of Bins to find given exploit for")
	exploitPtr := flag.String("exploit", "", "Exploit type:\n- bind-shell\n- capabilities\n- bin\n- file-download\n- file-read\n- file-upload\n- file-write\n- library-load\n- limited-suid\n- non-interactive-bind-shell\n- non-interactive-reverse-shell\n- reverse-shell\n- shell\n- sudo\n- suid")
	flag.Parse()

	return Flags{
		bins:    *binsPtr,
		exploit: *exploitPtr,
	}
}

func validateRequiredFlagValues(bins, exploit string) error {
	if bins == "" || exploit == "" {
		return fmt.Errorf("Error: make sure you set bins and an exploit")
	}

	if !isValidExploit(exploit) {
		return fmt.Errorf("Error: not a valid exploit")
	}

	return nil
}

func isValidExploit(exploit string) bool {
	validExploits := []string{
		"bind-shell",
		"capabilities",
		"bin",
		"file-download",
		"file-read",
		"file-upload",
		"file-write",
		"library-load",
		"limited-suid",
		"non-interactive-bind-shell",
		"non-interactive-reverse-shell",
		"reverse-shell",
		"shell",
		"sudo",
		"suid",
	}

	for _, valid := range validExploits {
		if exploit == valid {
			return true
		}
	}

	return false
}

func printFlagsBanner(exploit, bins string) {
	fmt.Print("\n")
	fmt.Print("---------------------------------")
	color.Note.Printf("\n EXPLOIT: %s", exploit)
	color.Note.Printf("\n BINS: %s\n", bins)
	fmt.Print("---------------------------------")
	fmt.Print("\n")
}

func printBins(url, bin, exploit string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch URL: %s, status: %s", url, resp.Status)
	}

	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return fmt.Errorf("error getting body: %w", err)
	}

	printContent(doc, bin, exploit, url)

	return nil
}

func printContent(doc *goquery.Document, bin, exploit, url string) {
	id := fmt.Sprintf("#%s", exploit)
	section := doc.Find(id)

	if section.Text() == "" {
		color.Danger.Printf("\n✘ %s not found\n", bin)
		return
	}

	printTitle(exploit, bin, url)
	printDescription(doc, id)
	printExamples(doc, id)
}

func printTitle(exploit, bin, url string) {
	exploitId := strings.ReplaceAll(exploit, " ", "-")
	hotpink.Printf("\n✔ %s %s/#%s\n", bin, url, exploitId)
}

func printDescription(doc *goquery.Document, id string) {
	text := ""
	doc.Find(id).NextUntil(".examples").Each(func(i int, s *goquery.Selection) {
		textTrimmed := strings.TrimSpace(s.Text())
		text += fmt.Sprintf("\n%s\n", textTrimmed)
	})

	fmt.Println(text)
}

func printExamples(doc *goquery.Document, id string) {
	codes := make([][]string, 0)

	doc.Find(id).NextFilteredUntil(".examples", "h2").Find("li").Each(func(i int, s *goquery.Selection) {
		s.Find("p").Each(func(i int, s *goquery.Selection) {
			liDescription := fmt.Sprintf("%s", strings.TrimSpace(s.Text()))
			codes = append(codes, []string{liDescription, ""})
		})

		s.Find("pre code").Each(func(i int, s *goquery.Selection) {
			code := fmt.Sprintf("%s\n", strings.TrimSpace(s.Text()))
			if len(codes) == i+1 {
				codes[i] = []string{codes[i][0], code}
			} else {
				codes = append(codes, []string{"", code})
			}
		})

	})

	for i, _ := range codes {
		if codes[i][0] != "" {
			fmt.Println(codes[i][0])
		}
		badass.Printf("%s\n", codes[i][1])
	}
}

func printCredits() {
	fmt.Print("\n")
	fmt.Print("--------------------------------------------------------------------------------------------")
	fmt.Print("\n\n")
	fmt.Print("- contribute to GTFOBins https://gtfobins.github.io/contribute/")
	fmt.Print("\n")
	fmt.Print("- follow GTFOBins' creators https://twitter.com/norbemi https://twitter.com/cyrus_and")
	fmt.Print("\n")
	fmt.Print("- follow me https://twitter.com/nightshiftc")
}
