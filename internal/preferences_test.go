package internal

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPrefs_SetBool(t *testing.T) {
	p := NewInMemoryPreferences()
	p.SetBool("testBool", true)

	assert.True(t, p.Bool("testBool"))
}

func TestPrefs_Bool(t *testing.T) {
	p := NewInMemoryPreferences()
	p.WriteValues(func(val map[string]any) {
		val["testBool"] = true
	})

	assert.True(t, p.Bool("testBool"))
	p.SetString("testBool", "fail")

	assert.False(t, p.Bool("testBool"))
}

func TestPrefs_BoolListWithFallback(t *testing.T) {
	p := NewInMemoryPreferences()
	boolstf := []bool{true, false}

	boolstt := []bool{true, true}

	assert.Equal(t, boolstf, p.BoolListWithFallback("testBoolList", boolstf))
	p.SetString("testBool", "fail")
	assert.Equal(t, boolstt, p.BoolListWithFallback("testBool", boolstt))
}

func TestPrefs_BoolWithFallback(t *testing.T) {
	p := NewInMemoryPreferences()

	assert.True(t, p.BoolWithFallback("testBool", true))
	p.SetBool("testBool", false)
	assert.Equal(t, 1, p.IntWithFallback("testBool", 1))
	p.SetString("testBool", "fail")
	assert.Equal(t, "fail", p.StringWithFallback("testBool", "fail"))
}

func TestPrefs_Bool_Zero(t *testing.T) {
	p := NewInMemoryPreferences()

	assert.False(t, p.Bool("testBool"))
}

func TestPrefs_SetFloat(t *testing.T) {
	p := NewInMemoryPreferences()
	p.SetFloat("testFloat", 1.7)

	assert.Equal(t, 1.7, p.Float("testFloat"))
}

func TestPrefs_Float(t *testing.T) {
	p := NewInMemoryPreferences()
	p.WriteValues(func(val map[string]any) {
		val["testFloat"] = 1.2
	})

	assert.Equal(t, 1.2, p.Float("testFloat"))
}

func TestPrefs_FloatWithFallback(t *testing.T) {
	p := NewInMemoryPreferences()

	assert.Equal(t, 1.0, p.FloatWithFallback("testFloat", 1.0))
	p.WriteValues(func(val map[string]any) {
		val["testFloat"] = 1.2
	})
	assert.Equal(t, 1.2, p.FloatWithFallback("testFloat", 1.0))

	assert.Equal(t, "bad", p.StringWithFallback("testFloat", "bad"))

	assert.Equal(t, 1.2, p.FloatWithFallback("testFloat", 1.3))

}

func TestPrefs_Float_Zero(t *testing.T) {
	p := NewInMemoryPreferences()

	assert.Equal(t, 0.0, p.Float("testFloat"))
}

func TestPrefs_SetInt(t *testing.T) {
	p := NewInMemoryPreferences()
	p.SetInt("testInt", 5)

	assert.Equal(t, 5, p.Int("testInt"))
}

func TestPrefs_Int(t *testing.T) {
	p := NewInMemoryPreferences()
	p.WriteValues(func(val map[string]any) {
		val["testInt"] = 5
	})
	assert.Equal(t, 5, p.Int("testInt"))
}

func TestPrefs_IntWithFallback(t *testing.T) {
	p := NewInMemoryPreferences()

	assert.Equal(t, 2, p.IntWithFallback("testInt", 2))
	p.WriteValues(func(val map[string]any) {
		val["testInt"] = 5
	})
	assert.Equal(t, 5, p.IntWithFallback("testInt", 2))

	assert.True(t, p.BoolWithFallback("testInt", true))

	assert.Equal(t, 5.0, p.FloatWithFallback("testInt", 1.2))

	assert.Equal(t, 5, p.IntWithFallback("testInt", 2))
}

func TestPrefs_Int_Zero(t *testing.T) {
	p := NewInMemoryPreferences()

	assert.Equal(t, 0, p.Int("testInt"))
}

func TestPrefs_SetString(t *testing.T) {
	p := NewInMemoryPreferences()
	p.SetString("test", "value")

	assert.Equal(t, "value", p.String("test"))
}

func TestPrefs_String(t *testing.T) {
	p := NewInMemoryPreferences()
	p.WriteValues(func(val map[string]any) {
		val["test"] = "value"
	})

	assert.Equal(t, "value", p.String("test"))
}

func TestPrefs_StringWithFallback(t *testing.T) {
	p := NewInMemoryPreferences()

	assert.Equal(t, "default", p.StringWithFallback("test", "default"))
	p.WriteValues(func(val map[string]any) {
		val["test"] = "value"
	})
	assert.Equal(t, "value", p.StringWithFallback("test", "default"))

	assert.True(t, p.BoolWithFallback("test", true))

	assert.Equal(t, "value", p.StringWithFallback("test", "default"))
}

func TestPrefs_String_Zero(t *testing.T) {
	p := NewInMemoryPreferences()

	assert.Equal(t, "", p.String("test"))
}

func TestInMemoryPreferences_OnChange(t *testing.T) {
	p := NewInMemoryPreferences()
	called := false
	p.AddChangeListener(func() {
		called = true
	})

	p.SetString("dummy", "another")
	time.Sleep(time.Millisecond * 100)
	assert.True(t, called)

	called = false
	p.RemoveValue("dummy")
	time.Sleep(time.Millisecond * 100)
	assert.True(t, called)
}

func TestRemoveValue(t *testing.T) {
	p := NewInMemoryPreferences()

	p.SetBool("dummy", true)
	p.SetFloat("pi", 3.14)
	p.SetInt("number", 2)
	p.SetString("month", "January")

	p.RemoveValue("dummy")
	p.RemoveValue("pi")
	p.RemoveValue("number")
	p.RemoveValue("month")

	assert.False(t, p.Bool("dummy"))
	assert.Equal(t, float64(0), p.Float("pi"))
	assert.Equal(t, 0, p.Int("number"))
	assert.Equal(t, "", p.String("month"))
}

func TestPrefs_SetSameValue(t *testing.T) {
	p := NewInMemoryPreferences()
	called := 0
	p.AddChangeListener(func() {
		called++
	})

	// We should not fire change when it hasn't changed.
	for i := 0; i < 2; i++ {
		p.SetBool("enabled", true)
		time.Sleep(time.Millisecond * 100)

		assert.Equal(t, 1, called)
	}

	p.SetBool("enabled", false)
	time.Sleep(time.Millisecond * 100)

	assert.Equal(t, 2, called)
}

func TestPrefs_SetSameSliceValue(t *testing.T) {
	p := NewInMemoryPreferences()
	called := 0
	p.AddChangeListener(func() {
		called++
	})

	// We should not fire change when it hasn't changed.
	for i := 0; i < 2; i++ {
		p.SetStringList("items", []string{"1", "2"})
		time.Sleep(time.Millisecond * 100)

		assert.Equal(t, 1, called)
	}

	p.SetStringList("items", []string{"3", "4"})
	time.Sleep(time.Millisecond * 100)

	assert.Equal(t, 2, called)
}
