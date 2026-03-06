package cmd

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/AkshayS96/bolt/internal/utils"

	"github.com/spf13/cobra"
)

func init() {
	timeCmd := &cobra.Command{
		Use:   "time",
		Short: "Time and timestamp utilities",
	}

	timeNowCmd := &cobra.Command{
		Use:   "now",
		Short: "Print current time in multiple formats",
		Run: func(cmd *cobra.Command, args []string) {
			now := time.Now()
			utils.PrintKeyValue("ISO 8601:", now.Format(time.RFC3339))
			utils.PrintKeyValue("Unix:", fmt.Sprintf("%d", now.Unix()))
			utils.PrintKeyValue("UTC:", now.UTC().Format(time.RFC3339))
			utils.PrintKeyValue("Local:", now.Format("2006-01-02 15:04:05 MST"))
		},
	}

	timeUnixCmd := &cobra.Command{
		Use:   "unix",
		Short: "Print current Unix timestamp",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(time.Now().Unix())
		},
	}

	timeISOCmd := &cobra.Command{
		Use:   "iso",
		Short: "Print current ISO 8601 timestamp",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(time.Now().Format(time.RFC3339))
		},
	}

	timeConvertCmd := &cobra.Command{
		Use:   "convert <timestamp>",
		Short: "Convert between Unix timestamp and ISO 8601",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			input := args[0]
			// Try parsing as Unix timestamp
			if unix, err := strconv.ParseInt(input, 10, 64); err == nil {
				t := time.Unix(unix, 0)
				utils.PrintKeyValue("ISO 8601:", t.Format(time.RFC3339))
				utils.PrintKeyValue("UTC:", t.UTC().Format(time.RFC3339))
				utils.PrintKeyValue("Local:", t.Format("2006-01-02 15:04:05 MST"))
				return
			}
			// Try parsing as ISO time
			formats := []string{
				time.RFC3339,
				"2006-01-02T15:04:05",
				"2006-01-02 15:04:05",
				"2006-01-02",
			}
			for _, f := range formats {
				if t, err := time.Parse(f, input); err == nil {
					utils.PrintKeyValue("Unix:", fmt.Sprintf("%d", t.Unix()))
					utils.PrintKeyValue("ISO 8601:", t.Format(time.RFC3339))
					return
				}
			}
			utils.PrintError("Could not parse timestamp: " + input)
		},
	}

	timeDiffCmd := &cobra.Command{
		Use:   "diff <date1> <date2>",
		Short: "Calculate difference between two dates",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			t1, err1 := parseTime(args[0])
			t2, err2 := parseTime(args[1])
			if err1 != nil {
				utils.PrintError("Cannot parse date1: " + err1.Error())
				return
			}
			if err2 != nil {
				utils.PrintError("Cannot parse date2: " + err2.Error())
				return
			}
			diff := t2.Sub(t1)
			absDiff := diff
			if absDiff < 0 {
				absDiff = -absDiff
			}
			days := int(math.Floor(absDiff.Hours() / 24))
			hours := int(math.Floor(absDiff.Hours())) % 24
			minutes := int(math.Floor(absDiff.Minutes())) % 60

			direction := "later"
			if diff < 0 {
				direction = "earlier"
			}

			utils.PrintKeyValue("Difference:", fmt.Sprintf("%d days, %d hours, %d minutes (%s)", days, hours, minutes, direction))
			utils.PrintKeyValue("Total hours:", fmt.Sprintf("%.1f", absDiff.Hours()))
			utils.PrintKeyValue("Total seconds:", fmt.Sprintf("%.0f", absDiff.Seconds()))
		},
	}

	timeCmd.AddCommand(timeNowCmd, timeUnixCmd, timeISOCmd, timeConvertCmd, timeDiffCmd)
	rootCmd.AddCommand(timeCmd)
}

func parseTime(s string) (time.Time, error) {
	// Try unix timestamp
	if unix, err := strconv.ParseInt(s, 10, 64); err == nil {
		return time.Unix(unix, 0), nil
	}
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unrecognized format: %s (use YYYY-MM-DD or RFC3339)", s)
}
