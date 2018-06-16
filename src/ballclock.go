package ballclock;

import(
    "fmt"
    "encoding/json"
    "github.com/jimlawless/whereami"
    "github.com/emirpasic/gods/stacks/arraystack"
    "github.com/emirpasic/gods/lists/arraylist"
)

type boolean interface {
    Compare(other interface{}) bool
}

type finite interface {
    Full() bool
}

type ball struct {
    int64
}

const(
    minLevelSlots = int64(5)
    fiveMinLevelSlots = int64(12)
    hourLevelSlots = int64(12)
)

func (b ball) Compare(other ball) bool {
    return b == other
}

func (b ball) toNative() int64 {
    return b.int64
}

type level struct {
    *arraystack.Stack
    max int64
}

type Clock struct {
    queue *arraylist.List
    min *level
    fivemin *level
    hour *level
    halfdays int64
}

func assert(cond bool, line interface{}) {
    if(!cond) {
        fmt.Printf("panic at: %v\n", line)
        panic("")
    }
}

func newBallList(nball int64) *arraylist.List {
    q := arraylist.New()
    for i := int64(0); i < nball; i++ {
        q.Add(ball{i})
    }
    return q
}

func newLevel(maxball int64) *level {
    return &level{Stack: arraystack.New(), max: maxball}
}

func (l *level) push(b ball) (ball, bool) {
    l.Push(b)
    if(l.Full()) {
        ba, _ := l.Pop()
        return ba.(ball), false
    } else {
        return ball{-1}, true
    }
}

func (l *level) Full() bool {
    return l.max == int64(l.Stack.Size())
}

func newClock(nball int64) *Clock {
    return &Clock{
        queue: newBallList(nball),
        min: newLevel(minLevelSlots),
        fivemin: newLevel(fiveMinLevelSlots),
        hour: newLevel(hourLevelSlots),
        halfdays: 0,
    }
}

func ballAdvance(l *level, r *level) (*ball, bool) {
    return nil, true
}

func (c *Clock) clockTick() bool {
    b, ok := c.queue.Get(0)
    assert(ok, whereami.WhereAmI())
    c.queue.Remove(0)

    var transfer = func (q *arraylist.List, s *level) {
        for !s.Empty() {
            b, ok := s.Pop()
            assert(ok, whereami.WhereAmI())
            if(ok) {
                q.Add(b.(ball))
            } else {
                assert(s.Size() == 0, whereami.WhereAmI())
                break;
            }
        }
    };

    b, ok = c.min.push(b.(ball))
    if(!ok) {
        transfer(c.queue, c.min)
        b, ok = c.fivemin.push(b.(ball))
        if(!ok) {
            transfer(c.queue, c.fivemin)
            b, ok = c.hour.push(b.(ball))
            if(!ok) {
                transfer(c.queue, c.hour)
                c.queue.Add(b.(ball))
                c.halfdays++
                return false
            }
        }
    }
    return true
}

func (c *Clock) print() {
    fmt.Println("queue:   ", c.queue)
    fmt.Println("min:     ", c.min)
    fmt.Println("fivemin: ", c.fivemin)
    fmt.Println("hour:    ", c.hour)
    fmt.Println("halfday: ", c.halfdays)
}

func compare(l *arraylist.List, r *arraylist.List) bool {
    if l.Size() != r.Size() {
        return false
    }
    for i := 0; i < l.Size(); i++ {
        lval, _ := l.Get(i)
        rval, _ := r.Get(i)
        if lval.(ball) != rval.(ball) {
            return false
        }
    }
    return true
}

func (c *Clock) Halfdays() int64 {
    return c.halfdays
}

func (c *Clock) Nball() int64 {
    return int64(
        c.min.Size() +
        c.fivemin.Size() +
        c.hour.Size() +
        c.queue.Size())
}

func (c *Clock) ToJson() []byte {
    var toSlice = func (arr []interface{}) []int64 {
        s := make([]int64, len(arr))
        for i := range arr {
            s[i] = arr[i].(ball).int64
        }
        return s
    }
    type intermediate struct {
        Min []int64
        FiveMin []int64
        Hour []int64
        Main []int64
    }
    b, err := json.Marshal(intermediate{
        Min: toSlice(c.min.Values()),
        FiveMin: toSlice(c.fivemin.Values()),
        Hour: toSlice(c.hour.Values()),
        Main: toSlice(c.queue.Values()),
    })
    assert(err == nil, whereami.WhereAmI())
    return b
}

func RunComplete(nball int64) *Clock {
    c := newClock(nball)
    target := newBallList(nball)
    for c.clockTick() || !compare(target, c.queue) {}
    return c
}

func RunMinutes(nball int64, min int64) *Clock {
    i := int64(0)
    c := newClock(nball)
    for i < min {
        c.clockTick()
        i++
    }
    return c
}
