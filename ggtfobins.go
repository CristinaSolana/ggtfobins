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
	commands string
	exploit string
}

func main() {
	flag.Usage = func() {
		color.Info.Print(" ______     ______     ______   ______   ______     ______     __     __   __     ______    \n/\\  ___\\   /\\  ___\\   /\\__  _\\ /\\  ___\\ /\\  __ \\   /\\  == \\   /\\ \\   /\\ \"-.\\ \\   /\\  ___\\   \n\\ \\ \\__ \\  \\ \\ \\__ \\  \\/_/\\ \\/ \\ \\  __\\ \\ \\ \\/\\ \\  \\ \\  __<   \\ \\ \\  \\ \\ \\-.  \\  \\ \\___  \\  \n \\ \\_____\\  \\ \\_____\\    \\ \\_\\  \\ \\_\\    \\ \\_____\\  \\ \\_____\\  \\ \\_\\  \\ \\_\\\\\"\\_\\  \\/\\_____\\ \n  \\/_____/   \\/_____/     \\/_/   \\/_/     \\/_____/   \\/_____/   \\/_/   \\/_/ \\/_/   \\/_____/")
		fmt.Print("\n\n")
		color.Info.Print("Get info about a given exploit for given commands\n")
		color.Info.Print("Usage: ggtfobins.go  --exploit suid --commands cpan,bash\n")
	}

	flags := getFlags()
	flag.Parse()

	commands := flags.commands
	exploit := flags.exploit
	err := validateRequiredFlagValues(commands, exploit)
	if err != nil {
		color.Red.Println(err)
		return
	}

	printExploitBanner(exploit)

	commandsList := strings.Split(commands, ",")
	for _, command := range commandsList {
		trimmedCommand := strings.TrimSpace(command)
		url := fmt.Sprintf("https://gtfobins.github.io/gtfobins/%s", trimmedCommand)
		err = printBins(url, trimmedCommand, exploit)
		if err != nil {
			log.Println(err)
		}
	}

	printCredits()
}

func getFlags () Flags {
	commandsPtr := flag.String("commands", "", "Comma-separated list of Commands to find given exploit for")
	exploitPtr := flag.String("exploit", "", "Exploit type:\n- bind-shell\n- capabilities\n- command\n- file-download\n- file-read\n- file-upload\n- file-write\n- library-load\n- limited-suid\n- non-interactive-bind-shell\n- non-interactive-reverse-shell\n- reverse-shell\n- shell\n- sudo\n- suid")
	flag.Parse()

	return Flags{
		commands: *commandsPtr,
		exploit: *exploitPtr,
	}
}

func validateRequiredFlagValues(commands, exploit string) error {
	if commands == "" {
		return fmt.Errorf("error: needz commands")
	}

	if exploit == "" {
		return fmt.Errorf("error: needz exploit type")
	}

	return nil
}

func printExploitBanner (exploit string) {
	fmt.Print("\n")
	fmt.Print("---------------------------------")
	color.Blue.Printf("\n EXPLOIT: %s\n", exploit)
	fmt.Print("---------------------------------")
	fmt.Print("\n")
}

func printBins(url, command, exploit string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	printContent(doc, command, exploit, url)

	return nil
}

func printContent(doc *goquery.Document, command, exploit, url string) {
	id := fmt.Sprintf("#%s", exploit)
	section := doc.Find(id)

	if section.Text() == "" {
		color.Yellow.Printf("\n✘ %s not found\n", command)
		return
	}

	printTitle(exploit, command, url)
	printDescription(doc, id)
	printExamples(doc, id)
}

func printTitle (exploit, command, url string) {
	mag := color.HEX("#f06")
	exploitId := strings.ReplaceAll(exploit, " ", "-")
	mag.Printf("\n✔ %s %s/#%s\n", command, url, exploitId)
}

func printDescription (doc *goquery.Document, id string) {
	text := ""
	doc.Find(id).NextUntil("ul").Each(func(i int, s *goquery.Selection) {
		textTrimmed := strings.TrimSpace(s.Text())
		text += fmt.Sprintf("%s\n", textTrimmed)
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
