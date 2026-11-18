package referenceframe

import (
	"math"
	"math/rand"
	"sync"
	"testing"

	"github.com/golang/geo/r3"

	"go.viam.com/rdk/spatialmath"
)

type unpreparedCustomFrame struct{ Frame }

func TestPreparedTransformRelativePrismaticAndFallback(t *testing.T) {
	fs := NewEmptyFrameSystem("prepared-test")
	rotation, err := NewRotationalFrame("rotation", spatialmath.R4AA{RZ: 1}, Limit{Min: -math.Pi, Max: math.Pi})
	if err != nil {
		t.Fatal(err)
	}
	translation, err := NewTranslationalFrame("slide", r3.Vector{X: 1}, Limit{Min: -200, Max: 200})
	if err != nil {
		t.Fatal(err)
	}
	tool, err := NewStaticFrame("tool", spatialmath.NewPoseFromPoint(r3.Vector{X: 20, Y: 30, Z: 40}))
	if err != nil {
		t.Fatal(err)
	}
	for _, pair := range [][2]Frame{{rotation, fs.World()}, {translation, rotation}, {tool, translation}} {
		if err := fs.AddFrame(pair[0], pair[1]); err != nil {
			t.Fatal(err)
		}
	}
	inputs := FrameSystemInputs{"rotation": {.3}, "slide": {10}, "tool": {}}.ToLinearInputs()
	schema, err := inputs.GetSchema(fs)
	if err != nil {
		t.Fatal(err)
	}
	for _, parent := range []string{World, "rotation", "slide", "tool"} {
		prepared, err := fs.PrepareTransform(schema, "tool", parent)
		if err != nil {
			t.Fatal(err)
		}
		rng := rand.New(rand.NewSource(23))
		for range 200 {
			inputs.Put("rotation", []Input{rng.Float64()*6 - 3})
			inputs.Put("slide", []Input{rng.Float64()*300 - 150})
			got, err := prepared.Evaluate(inputs.GetLinearizedInputs())
			if err != nil {
				t.Fatal(err)
			}
			want, err := fs.TransformToDQ(inputs, "tool", parent)
			if err != nil {
				t.Fatal(err)
			}
			if got.Number != want.Number {
				t.Fatalf("relative FK differs: %v versus %v", got.Number, want.Number)
			}
		}
		q := append([]float64(nil), inputs.GetLinearizedInputs()...)
		want, _ := prepared.Evaluate(q)
		var wg sync.WaitGroup
		for range 8 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for range 100 {
					got, err := prepared.Evaluate(q)
					if err != nil || got.Number != want.Number {
						t.Error("concurrent FK changed")
					}
				}
			}()
		}
		wg.Wait()
	}
	custom := unpreparedCustomFrame{NewZeroStaticFrame("custom")}
	if err := fs.AddFrame(custom, fs.World()); err != nil {
		t.Fatal(err)
	}
	if _, err := fs.PrepareTransform(schema, "custom", World); err == nil {
		t.Fatal("custom frames must retain native fallback")
	}
	if _, err := fs.TransformToDQ(inputs, "custom", World); err != nil {
		t.Fatal("native custom frame failed", err)
	}
	updated, _ := NewStaticFrame("tool", spatialmath.NewPoseFromPoint(r3.Vector{Z: 200}))
	if err := fs.ReplaceFrame(updated); err != nil {
		t.Fatal(err)
	}
	prepared, err := fs.PrepareTransform(schema, "tool", World)
	if err != nil {
		t.Fatal(err)
	}
	got, _ := prepared.Evaluate(inputs.GetLinearizedInputs())
	want, _ := fs.TransformToDQ(inputs, "tool", World)
	if got.Number != want.Number {
		t.Fatal("new snapshot did not incorporate tool change")
	}
}

func TestPreparedTransformMimicModel(t *testing.T) {
	model, err := ParseModelXMLFile("testfiles/test_mimic_serial.urdf", "arm", nil)
	if err != nil {
		t.Fatal(err)
	}
	fs := NewEmptyFrameSystem("mimic")
	if err := fs.AddFrame(model, fs.World()); err != nil {
		t.Fatal(err)
	}
	inputs := FrameSystemInputs{"arm": {.4, -.2}}.ToLinearInputs()
	schema, err := inputs.GetSchema(fs)
	if err != nil {
		t.Fatal(err)
	}
	prepared, err := fs.PrepareTransform(schema, "arm", World)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 100; i++ {
		inputs.Put("arm", []Input{float64(i)*.01 - .5, float64(i)*-.02 + 1})
		got, err := prepared.Evaluate(inputs.GetLinearizedInputs())
		if err != nil {
			t.Fatal(err)
		}
		want, err := fs.TransformToDQ(inputs, "arm", World)
		if err != nil {
			t.Fatal(err)
		}
		if got.Number != want.Number {
			t.Fatalf("mimic FK differs: %v vs %v", got.Number, want.Number)
		}
	}
}
