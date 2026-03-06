package cmd

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strconv"
	"strings"

	"github.com/AkshayS96/bolt/internal/utils"
	"github.com/disintegration/imaging"
	"github.com/spf13/cobra"
	_ "golang.org/x/image/webp"
)

func init() {
	var width, height int
	var ratio, bgColor string
	var blurSigma float64

	imgCmd := &cobra.Command{
		Use:   "img",
		Short: "Image processing utilities (resize, crop, info, optimize)",
	}

	// 1. img resize
	resizeCmd := &cobra.Command{
		Use:   "resize <input> <output>",
		Short: "Resize an image to a given width/height (maintains aspect ratio if one is 0)",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			input := args[0]
			output := args[1]
			if width <= 0 && height <= 0 {
				utils.PrintError("Must specify --width or --height")
				return
			}
			src, err := imaging.Open(input)
			if err != nil {
				utils.PrintError(fmt.Sprintf("Failed to open image: %v", err))
				return
			}
			dst := imaging.Resize(src, width, height, imaging.Lanczos)
			if err := imaging.Save(dst, output); err != nil {
				utils.PrintError(fmt.Sprintf("Failed to save image: %v", err))
				return
			}
			utils.PrintSuccess(fmt.Sprintf("Resized image saved to %s", output))
		},
	}
	resizeCmd.Flags().IntVarP(&width, "width", "w", 0, "Width to resize to")
	resizeCmd.Flags().IntVarP(&height, "height", "H", 0, "Height to resize to (0 to preserve aspect ratio)")

	// 2. img crop
	cropCmd := &cobra.Command{
		Use:   "crop <input> <output>",
		Short: "Crop an image to a specific aspect ratio (e.g. 16:9, 1:1)",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			input := args[0]
			output := args[1]
			parts := strings.Split(ratio, ":")
			if len(parts) != 2 {
				utils.PrintError("Ratio must be in the format W:H (e.g. 16:9)")
				return
			}
			rw, err1 := strconv.Atoi(parts[0])
			rh, err2 := strconv.Atoi(parts[1])
			if err1 != nil || err2 != nil || rw <= 0 || rh <= 0 {
				utils.PrintError("Invalid ratio values")
				return
			}

			src, err := imaging.Open(input)
			if err != nil {
				utils.PrintError(fmt.Sprintf("Failed to open image: %v", err))
				return
			}

			b := src.Bounds()
			srcW, srcH := b.Dx(), b.Dy()

			// Calculate target dimensions
			targetW := srcW
			targetH := (srcW * rh) / rw

			if targetH > srcH {
				targetH = srcH
				targetW = (srcH * rw) / rh
			}

			dst := imaging.Fill(src, targetW, targetH, imaging.Center, imaging.Lanczos)
			if err := imaging.Save(dst, output); err != nil {
				utils.PrintError(fmt.Sprintf("Failed to save image: %v", err))
				return
			}
			utils.PrintSuccess(fmt.Sprintf("Cropped image saved to %s", output))
		},
	}
	cropCmd.Flags().StringVarP(&ratio, "ratio", "r", "1:1", "Aspect ratio (e.g. 16:9, 4:3, 1:1)")

	// 3. img info
	infoCmd := &cobra.Command{
		Use:   "info <input>",
		Short: "Get basic information about an image",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			input := args[0]
			file, err := os.Open(input)
			if err != nil {
				utils.PrintError(fmt.Sprintf("Failed to open image: %v", err))
				return
			}
			defer file.Close()

			conf, format, err := image.DecodeConfig(file)
			if err != nil {
				utils.PrintError(fmt.Sprintf("Failed to decode image config: %v", err))
				return
			}

			stat, err := file.Stat()
			var sizeStr string
			if err == nil {
				sizeStr = fmt.Sprintf("%.2f KB", float64(stat.Size())/1024.0)
			} else {
				sizeStr = "Unknown"
			}

			utils.PrintKeyValue("Format:", strings.ToUpper(format))
			utils.PrintKeyValue("Width:", fmt.Sprintf("%d px", conf.Width))
			utils.PrintKeyValue("Height:", fmt.Sprintf("%d px", conf.Height))
			utils.PrintKeyValue("File Size:", sizeStr)
		},
	}

	// 4. img placeholder
	placeholderCmd := &cobra.Command{
		Use:   "placeholder <WxH> <output>",
		Short: "Generate a solid colored placeholder image",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			dim := args[0]
			output := args[1]

			parts := strings.Split(strings.ToLower(dim), "x")
			if len(parts) != 2 {
				utils.PrintError("Dimensions must be in the format WxH (e.g. 800x600)")
				return
			}
			w, err1 := strconv.Atoi(parts[0])
			h, err2 := strconv.Atoi(parts[1])
			if err1 != nil || err2 != nil || w <= 0 || h <= 0 {
				utils.PrintError("Invalid dimensions")
				return
			}

			hexColor := strings.TrimPrefix(bgColor, "#")
			var r, g, b uint8 = 200, 200, 200 // Default gray
			if len(hexColor) == 6 {
				parsed, err := strconv.ParseUint(hexColor, 16, 32)
				if err == nil {
					r = uint8(parsed >> 16)
					g = uint8((parsed >> 8) & 0xFF)
					b = uint8(parsed & 0xFF)
				}
			}

			dst := imaging.New(w, h, color.NRGBA{r, g, b, 255})
			if err := imaging.Save(dst, output); err != nil {
				utils.PrintError(fmt.Sprintf("Failed to save placeholder: %v", err))
				return
			}
			utils.PrintSuccess(fmt.Sprintf("Placeholder saved to %s", output))
		},
	}
	placeholderCmd.Flags().StringVarP(&bgColor, "color", "c", "cccccc", "Background color hex (e.g. ff5733)")

	// 5. img blur
	blurCmd := &cobra.Command{
		Use:   "blur <input> <output>",
		Short: "Apply a Gaussian blur to an image",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			input := args[0]
			output := args[1]

			if blurSigma <= 0 {
				utils.PrintError("Sigma must be > 0")
				return
			}

			src, err := imaging.Open(input)
			if err != nil {
				utils.PrintError(fmt.Sprintf("Failed to open image: %v", err))
				return
			}

			dst := imaging.Blur(src, blurSigma)
			if err := imaging.Save(dst, output); err != nil {
				utils.PrintError(fmt.Sprintf("Failed to save image: %v", err))
				return
			}
			utils.PrintSuccess(fmt.Sprintf("Blurred image saved to %s", output))
		},
	}
	blurCmd.Flags().Float64VarP(&blurSigma, "sigma", "s", 5.0, "Blur sigma (radius)")

	imgCmd.AddCommand(resizeCmd, cropCmd, infoCmd, placeholderCmd, blurCmd)
	rootCmd.AddCommand(imgCmd)
}
