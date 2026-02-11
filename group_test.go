package huh

import (
    "bytes"
    "fmt"
    "io"
    "sync/atomic"
    "testing"

    "github.com/charmbracelet/bubbles/key"
    tea "github.com/charmbracelet/bubbletea"
)

type testMsg struct {}

type Counters struct {
    numTestMsgs uint32
    numKeyMsgs uint32
    numOtherMsgs uint32
    numTotalMsgs uint32
}

type CountMsgsField struct {
    key string
    numTestMsgs uint32
    numKeyMsgs uint32
    numOtherMsgs uint32
    numTotalMsgs uint32
}

// A Field that just counts the number of tea.Msg it receives in Update
func NewCountMsgsField() *CountMsgsField {
    return &CountMsgsField{
        numTestMsgs: uint32(0),
        numKeyMsgs: uint32(0),
        numOtherMsgs: uint32(0),
        numTotalMsgs: uint32(0),
    }
}

func (_ *CountMsgsField) Init() tea.Cmd {
    return nil
}

func (c *CountMsgsField) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    atomic.AddUint32(&c.numTotalMsgs, 1)
    switch msg.(type) {
    case testMsg:
        atomic.AddUint32(&c.numTestMsgs, 1)
    case tea.KeyMsg:
        atomic.AddUint32(&c.numKeyMsgs, 1)
    default:
        atomic.AddUint32(&c.numOtherMsgs, 1)
    }
    return c, nil
}

func (_ *CountMsgsField) View() string {
    return ""
}

func (_ *CountMsgsField) Blur() tea.Cmd {
    return nil
}

func (_ *CountMsgsField) Focus() tea.Cmd {
    return nil
}

func (_ *CountMsgsField) Error() error {
    return nil
}

func (_ *CountMsgsField) Run() error {
    return nil
}

func (_ *CountMsgsField) RunAccessible(w io.Writer, r io.Reader) error {
    return nil
}

func (_ *CountMsgsField) Skip() bool {
    return false
}

func (_ *CountMsgsField) Zoom() bool {
    return false
}

func (_ *CountMsgsField) KeyBinds() []key.Binding {
    return []key.Binding{}
}

func (c *CountMsgsField) WithTheme(_ *Theme) Field {
    return c
}

func (c *CountMsgsField) WithKeyMap(_ *KeyMap) Field {
    return c
}

func (c *CountMsgsField) WithAccessible(_ bool) Field {
    return c
}

func (c *CountMsgsField) WithWidth(_ int) Field {
    return c
}

func (c *CountMsgsField) WithHeight(_ int) Field {
    return c
}

func (c *CountMsgsField) WithPosition(_ FieldPosition) Field {
    return c
}

func (c *CountMsgsField) GetKey() string {
    return c.key
}

func (c *CountMsgsField) GetValue() any {
    return Counters {
        numTestMsgs: atomic.LoadUint32(&c.numTestMsgs),
        numKeyMsgs: atomic.LoadUint32(&c.numKeyMsgs),
        numOtherMsgs: atomic.LoadUint32(&c.numOtherMsgs),
        numTotalMsgs: atomic.LoadUint32(&c.numTotalMsgs),
    }
}

func (c *CountMsgsField) Key(k string) *CountMsgsField {
    c.key = k
    return c
}

func testGroupSelectorMsgRedelivery[F interface { Field ; Key(string) F ; GetKey() string ; GetValue() any } ](t *testing.T, leaf F) (any, error) {
    numTestMsgs := uint32(0) // A Custom Msg (detect msg re-delivery bugs)
    numKeyMsgs := uint32(0)  // A Known-good Msg type (may detect test bugs)
    numOtherMsgs := uint32(0) // Noise Msgs (may detect test bugs)
    numTotalMsgs := uint32(0)

    // Input should be empty; Output is purely to keep
    // form rendering from obscuring 'go test' output.
    var in bytes.Buffer
    var out bytes.Buffer

    test_field := leaf.Key("Leaf")
    form := NewForm(NewGroup(test_field))

    countMsgs := func(_ tea.Model, msg tea.Msg) tea.Msg {
        atomic.AddUint32(&numTotalMsgs, 1)
        switch msg.(type) {
        case testMsg:
            atomic.AddUint32(&numTestMsgs, 1)
        case tea.KeyMsg:
            atomic.AddUint32(&numKeyMsgs, 1)
        default:
            atomic.AddUint32(&numOtherMsgs, 1)
        }
        return msg
    }

    p := tea.NewProgram(form,
        tea.WithInput(&in),
        tea.WithOutput(&out),
        tea.WithFilter(countMsgs),
    )

    finished := make(chan bool)

    // NOTE: The exact number of updateMsg sent is imprecise.
    go func() {
        p.Send(testMsg{})

        // Send tea.KeyMsg to test a past regression
        p.Send(tea.KeyMsg{Type: tea.KeyEnter}) // side effect: +tea.nextField

        // Need to do this to pump field value into form "results"
        // so that we can retrieve them via form.Get("Leaf") below.
        p.Send(form.NextGroup()) // side effect: +tea.nextField +tea.nextGroup

        p.Quit() // side effects: +tea.quitMsg
        p.Wait()
        finished <- true
    }()

    var model tea.Model
    var err error
    if model, err = p.Run(); err != nil {
        t.Fatal(err)
    }
    <-finished

    numTests := atomic.LoadUint32(&numTestMsgs)
    numKeys := atomic.LoadUint32(&numKeyMsgs)
    numOther := atomic.LoadUint32(&numOtherMsgs)
    totalMsgs := atomic.LoadUint32(&numTotalMsgs)

    if numTests != 1 {
        t.Errorf("Expected 1 testMsg deliveries to top level model got %d", numTests)
    }
    if numKeys != 1 {
        t.Errorf("Expected 1 tea.KeyMsg deliveries to top level model got %d", numKeys)
    }
    // account for side effects
    if numOther >= 4 {
        t.Logf("INFO: Estimated 4 or more other type of tea.Msg deliveries to top level model got %d", numOther)
    }

    expected := numTests + numKeys + numOther
    if totalMsgs != expected {
        // This is an error because the expected value is derived from the measured counts
        // rather than some static calculation.
        t.Errorf("Expected a total of %d tea.Msg deliveries to top level model got %d", expected, totalMsgs)
    }

    switch model.(type) {
        case *Form:
            form := model.(*Form)
            leaf_out := form.Get("Leaf")
            return leaf_out, nil
        default:
            return leaf, fmt.Errorf("tea returned a Model that is not a Form")
    }
}

// known-good Field type vs testing-only Field type (CountMsgsField)
func TestSelectMsgRedelivery(t *testing.T) {
    expected := "Selected A"
    s := NewSelect[string]().
            Options(NewOptions[string](expected, "Selected B")...)
    selected_any, err := testGroupSelectorMsgRedelivery(t, s)
    if err != nil {
        t.Fatalf("Count not retrieve selected %v", err)
    }
    selected := selected_any.(string)
    if selected != expected {
        t.Errorf("Expected \"%s\" got %s", expected, selected)
    }
}

func TestCountMsgRedelivery(t *testing.T) {
    initial_counter := NewCountMsgsField()
    counts_any, err := testGroupSelectorMsgRedelivery(t, initial_counter)
    if err != nil {
        t.Fatalf("Count not retrieve counters: %v", err)
    }
    counts := counts_any.(Counters)
    if counts.numTestMsgs != 1 {
        t.Errorf("Expected 1 testMsg deliveries to counter model got %d", counts.numTestMsgs)
    }
    if counts.numKeyMsgs != 1 {
        t.Errorf("Expected a 1 tea.KeyMsg deliveries to counter model got %d", counts.numKeyMsgs)
    }
    if counts.numOtherMsgs != 2 {
        t.Logf("INFO: Expected 2 other type of tea.Msg deliveries to counter model got %d", counts.numOtherMsgs)
    }

    expected := counts.numTestMsgs + counts.numKeyMsgs + counts.numOtherMsgs
    if counts.numTotalMsgs != expected {
        t.Errorf("Expected a total of %d tea.Msg deliveries to counter model got %d", expected, counts.numTotalMsgs)
    }
}
