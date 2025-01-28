package main

import (
    "bufio"
    "flag"
    "fmt"
    "image"
    "image/color"
    "image/png"
    "io"
    "math"
    "os"
    "path/filepath"
    "strings"

    "github.com/ajstarks/svgo"
)

const (
    pixelSize = 8
)

var colorMap = map[rune]color.RGBA{
    '0': {0, 0, 0, 255},         // black
    '1': {221, 0, 0, 255},       // red
    '2': {255, 206, 0, 255},     // yellow
}

func main() {
    decode := flag.Bool("d", false, "Decode PNG/SVG to values")
    blocksPerRow := flag.Int("b", 0, "Number of blocks per row (0 for single row)")
    useSVG := flag.Bool("v", false, "Use SVG format instead of PNG")
    help := flag.Bool("h", false, "Show help")
    flag.Parse()

    if *help || len(os.Args) == 1 {
        printUsage()
        os.Exit(0)
    }

    if *decode {
        if err := decodeToValues(os.Stdin, os.Stdout, *useSVG); err != nil {
            fmt.Fprintf(os.Stderr, "Error decoding: %v\n", err)
            os.Exit(1)
        }
    } else {
        if err := encodeValuesToImage(os.Stdin, os.Stdout, *blocksPerRow, *useSVG); err != nil {
            fmt.Fprintf(os.Stderr, "Error encoding: %v\n", err)
            os.Exit(1)
        }
    }
}

func printUsage() {
    fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
    fmt.Fprintln(os.Stderr, "  Encode: cat values.txt | "+filepath.Base(os.Args[0])+" -b blocks_per_row [-v] > output.png/svg")
    fmt.Fprintln(os.Stderr, "  Decode: cat input.png/svg | "+filepath.Base(os.Args[0])+" -d [-v] > output.txt")
    fmt.Fprintln(os.Stderr, "\nOptions:")
    flag.PrintDefaults()
}

func encodeValuesToImage(r io.Reader, w io.Writer, blocksPerRow int, useSVG bool) error {
    valueData, err := io.ReadAll(r)
    if err != nil {
        return fmt.Errorf("reading input: %w", err)
    }

    cleanValueData := strings.Map(func(r rune) rune {
        if r == ' ' || r == '\n' || r == '\r' {
            return -1
        }
        return r
    }, string(valueData))

    data := []byte(string(cleanValueData))

    blockCount := len(data)
    if blocksPerRow <= 0 {
        blocksPerRow = blockCount
    }

    rows := int(math.Ceil(float64(blockCount) / float64(blocksPerRow)))
    width := blocksPerRow * pixelSize
    height := rows * pixelSize

    if useSVG {
        return encodeSVG(w, data, width, height, blocksPerRow)
    }
    return encodePNG(w, data, width, height, blocksPerRow)
}

func encodePNG(w io.Writer, data []byte, width, height, blocksPerRow int) error {
    img := image.NewRGBA(image.Rect(0, 0, width, height))

    for i := 0; i < len(data); i++ {
        drawBlock(img, i, blocksPerRow, data[i])
    }

    return png.Encode(w, img)
}

func encodeSVG(w io.Writer, data []byte, width, height, blocksPerRow int) error {
    canvas := svg.New(w)
    canvas.Start(width, height)

    for i := 0; i < len(data); i++ {
        x, y := getBlockPosition(i, blocksPerRow)
        c := colorMap[rune(data[i])]
        canvas.Rect(x, y, pixelSize, pixelSize, fmt.Sprintf("fill:rgb(%d,%d,%d)", c.R, c.G, c.B))
    }

    canvas.End()
    return nil
}

func drawBlock(img *image.RGBA, blockIndex, blocksPerRow int, value byte) {
    x, y := getBlockPosition(blockIndex, blocksPerRow)
    c := colorMap[rune(value)]
    for dy := 0; dy < pixelSize; dy++ {
        for dx := 0; dx < pixelSize; dx++ {
            img.Set(x+dx, y+dy, c)
        }
    }
}

func getBlockPosition(blockIndex, blocksPerRow int) (x, y int) {
    return (blockIndex % blocksPerRow) * pixelSize, (blockIndex / blocksPerRow) * pixelSize
}

func decodeToValues(r io.Reader, w io.Writer, fromSVG bool) error {
    var data []byte
    var err error

    if fromSVG {
        data, err = decodeSVG(r)
    } else {
        data, err = decodePNG(r)
    }

    if err != nil {
        return err
    }

    // Remove padding
    for len(data) > 0 && data[len(data)-1] == 0 {
        data = data[:len(data)-1]
    }

    // Write value data
    _, err = fmt.Fprintf(w, "%s", string(data))
    if err != nil {
        return err
    }

    // Add a newline at the end
    _, err = fmt.Fprintln(w)
    return err
}

func decodePNG(r io.Reader) ([]byte, error) {
    img, err := png.Decode(r)
    if err != nil {
        return nil, fmt.Errorf("decoding PNG: %w", err)
    }

    bounds := img.Bounds()
    width, height := bounds.Max.X, bounds.Max.Y

    var data []byte

    for y := 0; y < height; y += pixelSize {
        for x := 0; x < width; x += pixelSize {
            r, g, b, _ := img.At(x, y).RGBA()
            c := color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), 255}
            for k, v := range colorMap {
                if v == c {
                    data = append(data, byte(k))
                    break
                }
            }
        }
    }

    return data, nil
}

func decodeSVG(r io.Reader) ([]byte, error) {
    var data []byte
    scanner := bufio.NewScanner(r)
    for scanner.Scan() {
        line := scanner.Text()
        if strings.Contains(line, "fill:rgb(") {
            colorStr := strings.Split(strings.Split(line, "fill:rgb(")[1], ")")[0]
            rgb := strings.Split(colorStr, ",")
            r := uint8(stringToInt(rgb[0]))
            g := uint8(stringToInt(rgb[1]))
            b := uint8(stringToInt(rgb[2]))
            c := color.RGBA{r, g, b, 255}
            for k, v := range colorMap {
                if v == c {
                    data = append(data, byte(k))
                    break
                }
            }
        }
    }
    return data, scanner.Err()
}

func stringToInt(s string) int {
    var result int
    fmt.Sscanf(s, "%d", &result)
    return result
}

