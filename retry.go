package confmysql

import (
	"time"

	"github.com/go-courier/envconf"
)

type Retry struct {
	Repeats  int
	Interval envconf.Duration
}

func (r *Retry) SetDefaults() {
	if r.Repeats == 0 {
		r.Repeats = 3
	}
	if r.Interval == 0 {
		r.Interval = envconf.Duration(10 * time.Second)
	}
}

func (r Retry) Do(exec func() error) (err error) {
	if r.Repeats <= 0 {
		err = exec()
		return
	}
	for i := 0; i < r.Repeats; i++ {
		err = exec()
		if err != nil {
			time.Sleep(time.Duration(r.Interval))
			continue
		}
		break
	}
	return
}
