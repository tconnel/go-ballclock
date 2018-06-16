package main;

import(
    "ballclock/src"
    "fmt"
    "strconv"
    "os"
    "time"
    "math"
    "github.com/docopt/docopt-go"
)

func float_div(n int64, d int64) float64 {
    return float64(n) / float64(d)
}

func mode1_output(nball int64, halfdays int64, elapsed time.Duration) {
    fmt.Printf("%v balls cycle after %v days.\n", nball, float_div(halfdays, 2))
    fmt.Printf("Completed in %v millseconds (%s seconds)\n",
            math.Round(float64(elapsed.Nanoseconds()) * float_div(int64(time.Nanosecond), int64(time.Millisecond))),
            strconv.FormatFloat(elapsed.Seconds(), 'f', 6, 64))
}

func mode2_output(c *ballclock.Clock) {
    fmt.Println(string(c.ToJson()))
}

func main() {
    argv := os.Args[1:]
    usage := `
Ball-Clock. Implements a Ball-Clock Simulation.

Usage:
  ballclock <nball>
  ballclock <nball> [--min=<minutes>]
  ballclock -h | --help
  ballclock --version

Options:
  -h --help           Show this screen.
  --version           Show version.
  --min=<minutes>     Minutes to run [default: 60 * 24 * 30 = 43200].`

    arguments, _ := docopt.ParseArgs(usage, argv, "0.0.1")
    i := arguments["<nball>"]
    min := arguments["--min"]
    nball, err_nball := strconv.ParseInt(i.(string), 0, 64)
    minVal, errMin := strconv.ParseInt(min.(string), 0, 64)
    if err_nball != nil || nball < int64(27) || nball > int64(127) {
        fmt.Println("<nball> must be integer in interval [27,127]")
        return
    }
    if errMin != nil {
        start := time.Now()
        c := ballclock.RunComplete(nball)
        elapsed := time.Since(start)
        mode1_output(nball, c.Halfdays(), elapsed)
    } else {
        mode2_output(ballclock.RunMinutes(nball, minVal))
    }
}
