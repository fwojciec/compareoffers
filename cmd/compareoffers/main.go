package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/fwojciec/compareoffers"
	"github.com/olekukonko/tablewriter"
)

const (
	defaultPrice            float64 = 38
	defaultPrintRuns        string  = "1000,2000,4000,8000,12000,20000,50000,100000"
	usage                   string  = "Usage: compareoffers [options] <offer> <offer>"
	offerPatternExplanation string  = "Offer pattern: ADVANCE__RATE-UNTIL_[...]_RATE (for example 1500__7-2000_8-4000_9)."
	offerRegex              string  = `^[0-9]{0,6}(\.[0-9]{1,2})?\__([0-9]{1,2}(\.[0-9]{1,2})?-[0-9]{1,6}_)*[0-9]{1,2}$`
	printRunsRegex          string  = `^([0-9]+,)+[0-9]+$`
)

func main() {
	m := NewMain()

	if err := m.ParseFlags(os.Args); err != nil {
		if err == flag.ErrHelp {
			os.Exit(0)
		}
		fmt.Fprintln(m.Stderr, err)
		os.Exit(1)
	}

	data := make([][]string, len(m.PrintRuns))
	for i, prun := range m.PrintRuns {
		o1 := m.Offers[0].CalcEarnings(m.Price, prun)
		o2 := m.Offers[1].CalcEarnings(m.Price, prun)
		d := o2 - o1
		row := []string{
			fmt.Sprintf("%d", prun),
			fmt.Sprintf("%.2f", o1),
			fmt.Sprintf("%.2f", o2),
			fmt.Sprintf("%.2f", d),
		}
		data[i] = row
	}
	table := tablewriter.NewWriter(m.Stdout)
	table.SetHeader([]string{"Sales level", "Offer 1", "Offer 2", "Difference"})
	table.AppendBulk(data)
	table.Render()
}

// Main represents the main program execution.
type Main struct {
	Stdout    io.Writer
	Stderr    io.Writer
	Price     float64
	PrintRuns []int
	Offers    []*compareoffers.Offer
}

// NewMain returns a new instance of Main
func NewMain() *Main {
	return &Main{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}

// ParseFlags parses CLI flags
func (m *Main) ParseFlags(args []string) error {
	fs := flag.NewFlagSet("compareoffers", flag.ContinueOnError)
	priceFlag := fs.Float64("price", defaultPrice, "price per copy")
	printrunsFlag := fs.String("printruns", defaultPrintRuns, "print runs")

	// custom usage
	fs.Usage = func() {
		fmt.Fprintf(m.Stdout, "%s\n\n%s\n\nOptions:\n", usage, offerPatternExplanation)
		fs.PrintDefaults()
	}

	if err := fs.Parse(args[1:]); err != nil {
		return err
	}

	// validate and set price
	if *priceFlag < 0 {
		return fmt.Errorf("price can't be negative")
	}
	m.Price = *priceFlag

	// parse, validate and set printruns
	prs, err := parsePrintRuns(*printrunsFlag)
	if err != nil {
		return err
	}
	m.PrintRuns = prs

	// parse, validate and set offers
	offerArgs := fs.Args()
	if len(offerArgs) != 2 {
		return fmt.Errorf("you must provide exactly two offers to compare")
	}
	offers := make([]*compareoffers.Offer, len(offerArgs))
	for i, toffer := range offerArgs {
		if !validateOfferInput(toffer) {
			return fmt.Errorf("invalid offer format")
		}
		parsed, err := compareoffers.NewOfferFromString(toffer)
		if err != nil {
			return err
		}
		offers[i] = parsed
	}
	m.Offers = offers
	return nil
}

func parsePrintRuns(raw string) ([]int, error) {
	raw = removeWhiteSpace(raw)
	r := regexp.MustCompile(printRunsRegex)
	if !r.MatchString(raw) {
		return nil, fmt.Errorf("invalid print runs format")
	}
	tps := strings.Split(raw, ",")
	pruns := make([]int, len(tps))
	for i, tp := range tps {
		// we know the fromat is valid so we can safely ignore the error
		p, _ := strconv.Atoi(tp)
		pruns[i] = p
	}
	return pruns, nil
}

func validateOfferInput(raw string) bool {
	r := regexp.MustCompile(offerRegex)
	return r.MatchString(raw)
}

func removeWhiteSpace(s string) string {
	space := regexp.MustCompile(`\s`)
	return space.ReplaceAllString(s, "")
}
