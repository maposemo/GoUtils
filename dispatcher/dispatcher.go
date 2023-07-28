package dispatcher

import (
	"fmt"
)

type Listener interface {
    Listen(event interface{})
}

type job struct {
    to string
    do interface{}
}

type Dispatcher struct {
    jobCh chan job
    table map[string]Listener
}

func NewDispatcher() *Dispatcher {
    d := &Dispatcher{
        jobCh: make(chan job, 10000),
        table: make(map[string]Listener),
    }

    go func() {
        for job := range d.jobCh {
            d.table[job.to].Listen(job.do)
        }
    }()

    return d
}

func (d *Dispatcher) Register(listener Listener, too ...string) error {
    for _, to := range too {
        if _, ok := d.table[to]; ok {
            return fmt.Errorf("Already %s address is used", to)
        }

        d.table[to] = listener
    }

    return nil
}

func (d *Dispatcher) Dispatch(to string, do interface{}) error {
    if _, ok := d.table[to]; !ok {
        return fmt.Errorf("'%s' address is not registered", to)
    }

    d.jobCh <- job{to: to, do: do}

    return nil
}
