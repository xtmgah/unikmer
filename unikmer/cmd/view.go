// Copyright © 2018 Wei Shen <shenwei356@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"io"
	"runtime"
	"strings"

	"github.com/shenwei356/unikmer"
	"github.com/shenwei356/xopen"
	"github.com/spf13/cobra"
)

// viewCmd represents
var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "read and output binary format to plain text",
	Long: `read and output binary format to plain text

`,
	Run: func(cmd *cobra.Command, args []string) {
		opt := getOptions(cmd)
		runtime.GOMAXPROCS(opt.NumCPUs)
		files := getFileList(args)

		outFile := getFlagString(cmd, "out-file")

		outfh, err := xopen.Wopen(outFile)
		checkError(err)
		defer outfh.Close()
		var infh *xopen.Reader
		var reader *unikmer.Reader
		var kcode unikmer.KmerCode

		for _, file := range files {
			if !isStdin(file) && !strings.HasSuffix(file, extDataFile) {
				checkError(fmt.Errorf("input should be stdin or %s file", extDataFile))
			}
			func() {
				infh, err = xopen.Ropen(file)
				checkError(err)
				defer infh.Close()

				reader, err = unikmer.NewReader(infh)
				checkError(err)

				for {
					kcode, err = reader.Read()
					if err != nil {
						if err == io.EOF {
							break
						}
						checkError(err)
					}

					// outfh.WriteString(fmt.Sprintf("%s\n", kcode.Bytes())) // slower
					outfh.WriteString(kcode.String() + "\n")
				}

			}()
		}
	},
}

func init() {
	RootCmd.AddCommand(viewCmd)

	viewCmd.Flags().StringP("out-file", "o", "-", `out file ("-" for stdout, suffix .gz for gzipped out)`)

}
