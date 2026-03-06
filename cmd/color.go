package cmd

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"bolt/internal/utils"

	"github.com/spf13/cobra"
)

func init() {
	colorCmd := &cobra.Command{
		Use:   "color",
		Short: "Color conversion tools",
	}

	colorHex2RGBCmd := &cobra.Command{
		Use:   "hex2rgb <hex>",
		Short: "Convert hex color to RGB",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			hex := strings.TrimPrefix(args[0], "#")
			if len(hex) == 3 {
				hex = string(hex[0]) + string(hex[0]) + string(hex[1]) + string(hex[1]) + string(hex[2]) + string(hex[2])
			}
			if len(hex) != 6 {
				utils.PrintError("Invalid hex color. Use format: #RRGGBB or RRGGBB")
				return
			}
			r, err1 := strconv.ParseInt(hex[0:2], 16, 64)
			g, err2 := strconv.ParseInt(hex[2:4], 16, 64)
			b, err3 := strconv.ParseInt(hex[4:6], 16, 64)
			if err1 != nil || err2 != nil || err3 != nil {
				utils.PrintError("Invalid hex color value")
				return
			}
			utils.PrintKeyValue("Hex:", "#"+strings.ToUpper(hex))
			utils.PrintKeyValue("RGB:", fmt.Sprintf("rgb(%d, %d, %d)", r, g, b))
			utils.PrintKeyValue("Values:", fmt.Sprintf("R=%d G=%d B=%d", r, g, b))
		},
	}

	colorRGB2HexCmd := &cobra.Command{
		Use:   "rgb2hex <r> <g> <b>",
		Short: "Convert RGB values to hex color",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			r, err1 := strconv.Atoi(args[0])
			g, err2 := strconv.Atoi(args[1])
			b, err3 := strconv.Atoi(args[2])
			if err1 != nil || err2 != nil || err3 != nil || r < 0 || r > 255 || g < 0 || g > 255 || b < 0 || b > 255 {
				utils.PrintError("Invalid RGB values. Each must be 0-255")
				return
			}
			hex := fmt.Sprintf("#%02X%02X%02X", r, g, b)
			utils.PrintKeyValue("RGB:", fmt.Sprintf("rgb(%d, %d, %d)", r, g, b))
			utils.PrintKeyValue("Hex:", hex)
		},
	}

	colorRandomCmd := &cobra.Command{
		Use:   "random",
		Short: "Generate a random color",
		Run: func(cmd *cobra.Command, args []string) {
			r := rand.Intn(256)
			g := rand.Intn(256)
			b := rand.Intn(256)
			hex := fmt.Sprintf("#%02X%02X%02X", r, g, b)
			utils.PrintKeyValue("Hex:", hex)
			utils.PrintKeyValue("RGB:", fmt.Sprintf("rgb(%d, %d, %d)", r, g, b))
			utils.PrintKeyValue("HSL:", rgbToHSL(r, g, b))
		},
	}

	colorCmd.AddCommand(colorHex2RGBCmd, colorRGB2HexCmd, colorRandomCmd)
	rootCmd.AddCommand(colorCmd)
}

func rgbToHSL(r, g, b int) string {
	rf := float64(r) / 255.0
	gf := float64(g) / 255.0
	bf := float64(b) / 255.0

	max := rf
	if gf > max {
		max = gf
	}
	if bf > max {
		max = bf
	}
	min := rf
	if gf < min {
		min = gf
	}
	if bf < min {
		min = bf
	}

	l := (max + min) / 2.0
	if max == min {
		return fmt.Sprintf("hsl(0, 0%%, %.0f%%)", l*100)
	}

	d := max - min
	s := d / (1 - abs(2*l-1))

	var h float64
	switch max {
	case rf:
		h = (gf - bf) / d
		if gf < bf {
			h += 6
		}
	case gf:
		h = (bf-rf)/d + 2
	case bf:
		h = (rf-gf)/d + 4
	}
	h *= 60

	return fmt.Sprintf("hsl(%.0f, %.0f%%, %.0f%%)", h, s*100, l*100)
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
