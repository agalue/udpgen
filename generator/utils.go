package generator

import (
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/valyala/fasttemplate"
)

func parseText(template string) string {
	t := fasttemplate.New(template, "[", "]")
	return t.ExecuteFuncString(func(w io.Writer, tag string) (int, error) {
		if strings.HasPrefix(tag, "rand_n:") {
			vars := strings.Split(tag, ":")
			if len(vars) == 2 {
				if max, err := strconv.Atoi(vars[1]); err == nil {
					return w.Write([]byte(strconv.Itoa(rand.Intn(max))))
				}
			}
		}
		if tag == "rand_ipv4" {
			ip := fmt.Sprintf("%d.%d.%d.%d", rand.Intn(255), rand.Intn(255), rand.Intn(255), rand.Intn(255))
			return w.Write([]byte(ip))
		}
		return 0, nil
	})
}

func init() {
	rand.Seed(time.Now().Unix())
}
