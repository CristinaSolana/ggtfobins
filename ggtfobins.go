package main

import (
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gookit/color"
	"log"
	"net/http"
	"strings"
)

type Flags struct {
	bins string
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

func printBanner () {
	color.Note.Print(" ______     ______     ______   ______   ______     ______     __     __   __     ______    \n/\\  ___\\   /\\  ___\\   /\\__  _\\ /\\  ___\\ /\\  __ \\   /\\  == \\   /\\ \\   /\\ \"-.\\ \\   /\\  ___\\   \n\\ \\ \\__ \\  \\ \\ \\__ \\  \\/_/\\ \\/ \\ \\  __\\ \\ \\ \\/\\ \\  \\ \\  __<   \\ \\ \\  \\ \\ \\-.  \\  \\ \\___  \\  \n \\ \\_____\\  \\ \\_____\\    \\ \\_\\  \\ \\_\\    \\ \\_____\\  \\ \\_____\\  \\ \\_\\  \\ \\_\\\\\"\\_\\  \\/\\_____\\ \n  \\/_____/   \\/_____/     \\/_/   \\/_/     \\/_____/   \\/_____/   \\/_/   \\/_/ \\/_/   \\/_____/")
	fmt.Print("\n\n")
	color.Note.Print("Get info about a given exploit for given bins\n")
	color.Note.Print("Usage: ggtfobins.go  --exploit suid --bins cpan,bash\n\n")
}

func getFlags () Flags {
	binsPtr := flag.String("bins", "", "Comma-separated list of Bins to find given exploit for")
	exploitPtr := flag.String("exploit", "", "Exploit type:\n- bind-shell\n- capabilities\n- bin\n- file-download\n- file-read\n- file-upload\n- file-write\n- library-load\n- limited-suid\n- non-interactive-bind-shell\n- non-interactive-reverse-shell\n- reverse-shell\n- shell\n- sudo\n- suid")
	flag.Parse()

	return Flags{
		bins: *binsPtr,
		exploit: *exploitPtr,
	}
}

func validateRequiredFlagValues(bins, exploit string) error {
	if bins == "" && exploit == "" {
		return fmt.Errorf("error: missing bins and exploit")
	}

	if bins == "" {
		return fmt.Errorf("error: missing bins")
	}

	if exploit == "" {
		return fmt.Errorf("error: missing exploit type")
	}

	return nil
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
		return err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
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

func printTitle (exploit, bin, url string) {
	mag := color.HEX("#f06")
	exploitId := strings.ReplaceAll(exploit, " ", "-")
	mag.Printf("\n✔ %s %s/#%s\n", bin, url, exploitId)
}

func printDescription (doc *goquery.Document, id string) {
	text := ""
	doc.Find(id).NextUntil(".examples").Each(func(i int, s *goquery.Selection) {
		textTrimmed := strings.TrimSpace(s.Text())
		text += fmt.Sprintf("\n%s\n", textTrimmed)
	})

	fmt.Println(text)
}

func printExamples (doc *goquery.Document, id string) {
	codes := make([][]string, 0)

	doc.Find(id).NextFilteredUntil(".examples", "h2").Find("li p").Each(func(i int, s *goquery.Selection) {
		liDescription := fmt.Sprintf("%s", strings.TrimSpace(s.Text()))
		codes = append(codes, []string{liDescription, ""})
	})

	doc.Find(id).NextFilteredUntil(".examples", "h2").Find("li pre code").Each(func(i int, s *goquery.Selection) {
		code := fmt.Sprintf("%s\n", strings.TrimSpace(s.Text()))
		if len(codes) == i + 1 {
			codes[i] = []string{codes[i][0], code}
		} else {
			codes = append(codes, []string{"", code})
		}
	})

	badass := color.HEX("#bada55")
	for i, _ := range codes {
		if codes[i][0] != "" {
			fmt.Println(codes[i][0])
		}
		badass.Printf("%s\n", codes[i][1])
	}
}

func printCredits () {
	fmt.Print("\n")
	fmt.Print("--------------------------------------------------------------------------------------------")
	fmt.Print("\n\n")
	fmt.Print("- contribute to GTFOBins https://gtfobins.github.io/contribute/")
	fmt.Print("\n")
	fmt.Print("- follow GTFOBins' creators https://twitter.com/norbemi https://twitter.com/cyrus_and")
	fmt.Print("\n")
	fmt.Print("- follow me https://twitter.com/nightshiftc")
}
