package balancer

import (
	"github.com/zehuamama/balancer/utils"
	"log"
	"net/url"
	"strconv"
	"sync"
)

type WeightRound struct {
	hosts         map[string]int
	currentWeight map[string]int
	initWeight    map[string]int
	sync.RWMutex
}

func NewWeightRoundRobin(_ []string, args any) Balancer {
	b := &WeightRound{
		hosts:         map[string]int{},
		currentWeight: map[string]int{},
		initWeight:    map[string]int{},
	}
	if args == nil {
		return b
	}
	for _, host := range args.([]string) {
		l := utils.SplitStringBySpaces(host)
		host, err := url.Parse(l[0])
		if err != nil {
			log.Println(err)
			continue
		}
		h := utils.GetHost(host)
		if len(l) == 2 {
			w, err := strconv.Atoi(l[1])
			if err != nil || w <= 0 {
				w = 1
			}
			b.hosts[h] = w
			b.currentWeight[h] = w
			b.initWeight[h] = w
		}
		if len(l) == 1 {
			b.hosts[h] = 1
			b.currentWeight[h] = 1
			b.initWeight[h] = 1
		}

	}
	return b
}

func init() {
	factories[WeightRoundBalancer] = NewWeightRoundRobin
}

func (r *WeightRound) Add(host string) {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.initWeight[host]; !ok {
		return
	}
	for h, _ := range r.hosts {
		if h == host {
			return
		}
	}
	r.hosts[host] = r.initWeight[host]
	r.currentWeight[host] = r.initWeight[host]
}

func (r *WeightRound) Remove(host string) {
	r.Lock()
	defer r.Unlock()
	for h, _ := range r.hosts {
		if h == host {
			delete(r.hosts, h)
			delete(r.currentWeight, h)
		}
	}
}

func (r *WeightRound) Balance(_ string) (string, error) {
	r.RLock()
	defer r.RUnlock()
	sumw := r.sumWeight()
	if sumw == 0 || len(r.hosts) == 0 {
		return "", NoHostError
	}
	var ret string
	wt := 0
	for host, _ := range r.currentWeight {
		r.currentWeight[host] += r.hosts[host]
	}
	for host, w := range r.currentWeight {
		if w >= wt {
			ret = host
			wt = w
		}
	}
	r.currentWeight[ret] -= sumw
	return ret, nil
}

func (r *WeightRound) Inc(host string) {}

func (r *WeightRound) Done(host string) {}

func (r *WeightRound) sumWeight() int {
	sum := 0
	for _, w := range r.hosts {
		sum += w
	}
	return sum
}
