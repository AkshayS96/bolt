package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/AkshayS96/bolt/internal/utils"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	cronCmd := &cobra.Command{
		Use:   "cron [expression]",
		Short: "Cron expression tools",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 1 {
				explainCron(args[0])
			} else {
				cmd.Help()
			}
		},
	}

	cronExplainCmd := &cobra.Command{
		Use:   "explain <expression>",
		Short: "Explain a cron expression in human-readable format",
		Long: `Explain a cron expression in human-readable format.
  
  Example: bolt cron explain "*/5 * * * *"
  
  Cron format: minute hour day-of-month month day-of-week`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			explainCron(args[0])
		},
	}

	cronCmd.AddCommand(cronExplainCmd)

	// --- cron build ---
	var cronMinute, cronHour, cronDay, cronMonth, cronWeekday string
	var cronEvery string

	cronBuildCmd := &cobra.Command{
		Use:   "build",
		Short: "Build a cron expression interactively",
		Long: `Build a cron expression using flags.

  Examples:
    bolt cron build --at 9:30                     → 30 9 * * *
    bolt cron build --every 5m                    → */5 * * * *
    bolt cron build --every 2h                    → 0 */2 * * *
    bolt cron build --at 9:00 --on mon-fri        → 0 9 * * 1-5
    bolt cron build --minute 0 --hour 9 --day 1   → 0 9 1 * *`,
		Run: func(cmd *cobra.Command, args []string) {
			min, hour, day, month, wday := cronMinute, cronHour, cronDay, cronMonth, cronWeekday

			// --every shortcut
			if cronEvery != "" {
				e := cronEvery
				if strings.HasSuffix(e, "m") {
					val := strings.TrimSuffix(e, "m")
					min = "*/" + val
					hour = "*"
				} else if strings.HasSuffix(e, "h") {
					val := strings.TrimSuffix(e, "h")
					min = "0"
					hour = "*/" + val
				} else if strings.HasSuffix(e, "d") {
					min = "0"
					hour = "0"
				} else {
					min = "*/" + e
				}
			}

			// --at shortcut (e.g. "9:30" or "14:00")
			atVal, _ := cmd.Flags().GetString("at")
			if atVal != "" {
				timeParts := strings.Split(atVal, ":")
				if len(timeParts) == 2 {
					min = timeParts[1]
					hour = timeParts[0]
				} else if len(timeParts) == 1 {
					min = "0"
					hour = timeParts[0]
				}
			}

			// --on shortcut (e.g. "mon-fri", "1,3,5")
			onVal, _ := cmd.Flags().GetString("on")
			if onVal != "" {
				wday = parseDayNames(onVal)
			}

			expr := fmt.Sprintf("%s %s %s %s %s", min, hour, day, month, wday)
			fmt.Println()
			explainCron(expr)
			fmt.Println()
			bold := color.New(color.Bold)
			bold.Print("  Copy: ")
			fmt.Println(expr)
			fmt.Println()
		},
	}

	cronBuildCmd.Flags().StringVar(&cronMinute, "minute", "*", "Minute field (0-59)")
	cronBuildCmd.Flags().StringVar(&cronHour, "hour", "*", "Hour field (0-23)")
	cronBuildCmd.Flags().StringVar(&cronDay, "day", "*", "Day of month (1-31)")
	cronBuildCmd.Flags().StringVar(&cronMonth, "month", "*", "Month (1-12)")
	cronBuildCmd.Flags().StringVar(&cronWeekday, "weekday", "*", "Day of week (0-7)")
	cronBuildCmd.Flags().StringVar(&cronEvery, "every", "", "Shortcut: 5m, 2h, 1d")
	cronBuildCmd.Flags().String("at", "", "Shortcut: time like 9:30 or 14:00")
	cronBuildCmd.Flags().String("on", "", "Shortcut: weekdays like mon-fri or mon,wed,fri")

	cronCmd.AddCommand(cronBuildCmd)
	rootCmd.AddCommand(cronCmd)
}

func explainCron(expr string) {
	parts := strings.Fields(expr)
	if len(parts) < 5 || len(parts) > 6 {
		utils.PrintError("Invalid cron expression. Expected 5 or 6 fields: minute hour day month weekday [year]")
		return
	}

	// Validate each field
	fieldDefs := []struct {
		name string
		min  int
		max  int
	}{
		{"Minute", 0, 59},
		{"Hour", 0, 23},
		{"Day of month", 1, 31},
		{"Month", 1, 12},
		{"Day of week", 0, 7},
	}
	for i, def := range fieldDefs {
		if err := validateCronField(parts[i], def.name, def.min, def.max); err != nil {
			utils.PrintError(err.Error())
			return
		}
	}

	accent := color.New(color.FgCyan, color.Bold)
	boldC := color.New(color.Bold)
	dimC := color.New(color.Faint)

	// Natural language first — the main output
	accent.Print("  ⏰ ")
	boldC.Println(cronToHuman(parts))
	fmt.Println()

	// Show the expression + field breakdown as secondary info
	dimC.Println("  Expression: " + expr)
	dimC.Printf("  Fields:     %-5s%-5s%-5s%-5s%s\n",
		parts[0], parts[1], parts[2], parts[3], parts[4])
	dimC.Printf("              min  hr   day  mon  wday\n")
}

func validateCronField(field, name string, min, max int) error {
	if field == "*" {
		return nil
	}

	// Step values like */5
	if strings.HasPrefix(field, "*/") {
		val := strings.TrimPrefix(field, "*/")
		n, err := strconv.Atoi(val)
		if err != nil || n < 1 || n > max {
			return fmt.Errorf("invalid %s: '%s' — step value must be between 1 and %d", name, field, max)
		}
		return nil
	}

	// Reject bare /N (not valid cron)
	if strings.HasPrefix(field, "/") {
		return fmt.Errorf("invalid %s: '%s' — did you mean '*%s'?", name, field, field)
	}

	// Ranges like 1-5
	if strings.Contains(field, "-") {
		rangeParts := strings.Split(field, "-")
		if len(rangeParts) != 2 {
			return fmt.Errorf("invalid %s: '%s' — bad range format", name, field)
		}
		from, err1 := strconv.Atoi(rangeParts[0])
		to, err2 := strconv.Atoi(rangeParts[1])
		if err1 != nil || err2 != nil || from < min || to > max || from > to {
			return fmt.Errorf("invalid %s: '%s' — range must be between %d and %d", name, field, min, max)
		}
		return nil
	}

	// Comma-separated like 1,3,5
	if strings.Contains(field, ",") {
		for _, v := range strings.Split(field, ",") {
			n, err := strconv.Atoi(strings.TrimSpace(v))
			if err != nil || n < min || n > max {
				return fmt.Errorf("invalid %s: '%s' — values must be between %d and %d", name, field, min, max)
			}
		}
		return nil
	}

	// Single numeric value
	n, err := strconv.Atoi(field)
	if err != nil {
		return fmt.Errorf("invalid %s: '%s' — expected a number, *, or pattern like */5", name, field)
	}
	if n < min || n > max {
		return fmt.Errorf("invalid %s: '%s' — value must be between %d and %d", name, field, min, max)
	}
	return nil
}

func parseDayNames(input string) string {
	dayMap := map[string]string{
		"sun": "0", "mon": "1", "tue": "2", "wed": "3",
		"thu": "4", "fri": "5", "sat": "6",
		"sunday": "0", "monday": "1", "tuesday": "2", "wednesday": "3",
		"thursday": "4", "friday": "5", "saturday": "6",
	}
	input = strings.ToLower(input)

	// Handle ranges like mon-fri
	if strings.Contains(input, "-") {
		parts := strings.Split(input, "-")
		if len(parts) == 2 {
			from, ok1 := dayMap[parts[0]]
			to, ok2 := dayMap[parts[1]]
			if ok1 && ok2 {
				return from + "-" + to
			}
		}
		return input
	}

	// Handle comma-separated like mon,wed,fri
	if strings.Contains(input, ",") {
		parts := strings.Split(input, ",")
		var nums []string
		for _, p := range parts {
			if num, ok := dayMap[strings.TrimSpace(p)]; ok {
				nums = append(nums, num)
			} else {
				nums = append(nums, strings.TrimSpace(p))
			}
		}
		return strings.Join(nums, ",")
	}

	// Single day
	if num, ok := dayMap[input]; ok {
		return num
	}
	return input
}
func explainField(field, label string) string {
	if field == "*" {
		return "every " + strings.ToLower(label)
	}
	if strings.HasPrefix(field, "*/") {
		interval := strings.TrimPrefix(field, "*/")
		return fmt.Sprintf("every %s %s(s)", interval, strings.ToLower(label))
	}
	if strings.Contains(field, ",") {
		return "at " + strings.ToLower(label) + "(s) " + field
	}
	if strings.Contains(field, "-") {
		parts := strings.Split(field, "-")
		return fmt.Sprintf("from %s to %s", parts[0], parts[1])
	}
	return field
}

func cronToHuman(parts []string) string {
	days := map[string]string{
		"0": "Sunday", "1": "Monday", "2": "Tuesday", "3": "Wednesday",
		"4": "Thursday", "5": "Friday", "6": "Saturday", "7": "Sunday",
	}
	months := map[string]string{
		"1": "January", "2": "February", "3": "March", "4": "April",
		"5": "May", "6": "June", "7": "July", "8": "August",
		"9": "September", "10": "October", "11": "November", "12": "December",
	}

	min, hour, dom, mon, dow := parts[0], parts[1], parts[2], parts[3], parts[4]

	// Build the "when" part (minute + hour)
	timeStr := ""
	switch {
	case min == "*" && hour == "*":
		timeStr = "Runs every minute"
	case (min == "*/1" || strings.HasPrefix(min, "*/")) && hour == "*":
		n := strings.TrimPrefix(min, "*/")
		if n == "1" {
			timeStr = "Runs every minute"
		} else {
			timeStr = fmt.Sprintf("Runs every %s minutes", n)
		}
	case min == "*" && strings.HasPrefix(hour, "*/"):
		n := strings.TrimPrefix(hour, "*/")
		if n == "1" {
			timeStr = "Runs every hour"
		} else {
			timeStr = fmt.Sprintf("Runs every %s hours", n)
		}
	case strings.HasPrefix(hour, "*/"):
		n := strings.TrimPrefix(hour, "*/")
		if n == "1" {
			timeStr = fmt.Sprintf("Runs every hour at minute %s", min)
		} else {
			timeStr = fmt.Sprintf("Runs every %s hours at minute %s", n, min)
		}
	case hour == "*":
		timeStr = fmt.Sprintf("Runs every hour at minute %s", min)
	default:
		timeStr = fmt.Sprintf("Runs at %s", formatTime(hour, min))
	}

	// Build the "date" parts
	var dateParts []string

	// Day of week
	if dow != "*" {
		if name, ok := days[dow]; ok {
			dateParts = append(dateParts, "every "+name)
		} else if strings.Contains(dow, "-") {
			rangeParts := strings.Split(dow, "-")
			from, to := rangeParts[0], rangeParts[1]
			if d, ok := days[from]; ok {
				from = d
			}
			if d, ok := days[to]; ok {
				to = d
			}
			dateParts = append(dateParts, from+" through "+to)
		} else if strings.Contains(dow, ",") {
			nums := strings.Split(dow, ",")
			var names []string
			for _, n := range nums {
				if d, ok := days[strings.TrimSpace(n)]; ok {
					names = append(names, d)
				}
			}
			if len(names) > 0 {
				dateParts = append(dateParts, "on "+strings.Join(names, ", "))
			}
		}
	}

	// Day of month
	if dom != "*" {
		dateParts = append(dateParts, "on day "+dom+" of the month")
	}

	// Month
	if mon != "*" {
		if name, ok := months[mon]; ok {
			dateParts = append(dateParts, "in "+name)
		} else {
			dateParts = append(dateParts, "in month "+mon)
		}
	}

	if len(dateParts) > 0 {
		return timeStr + ", " + strings.Join(dateParts, ", ")
	}
	return timeStr
}

func formatTime(hour, min string) string {
	h, err := strconv.Atoi(hour)
	if err != nil {
		return hour + ":" + min
	}
	m, err := strconv.Atoi(min)
	if err != nil {
		return hour + ":" + min
	}
	period := "AM"
	display := h
	if h >= 12 {
		period = "PM"
		if h > 12 {
			display = h - 12
		}
	}
	if h == 0 {
		display = 12
	}
	return fmt.Sprintf("%d:%02d %s", display, m, period)
}
