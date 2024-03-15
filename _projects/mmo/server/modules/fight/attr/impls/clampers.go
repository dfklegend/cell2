package impls

import (
	"mmo/modules/fight/attr"
)

// MinPercentClapmper 最小衰减到多少
type MinPercentClapmper struct {
	min float64
}

func NewMinPercentClapmper(min float64) *MinPercentClapmper {
	return &MinPercentClapmper{
		min: min,
	}
}

func (c *MinPercentClapmper) Clamp(v float64) float64 {
	if v < c.min {
		return c.min
	}
	return v
}

// ----

// MinValueClapmper 最小值
type MinValueClapmper struct {
	min attr.Value
}

func NewMinValueClapmper(min attr.Value) *MinValueClapmper {
	return &MinValueClapmper{
		min: min,
	}
}

func (c *MinValueClapmper) Clamp(v attr.Value) attr.Value {
	if v < c.min {
		return c.min
	}
	return v
}

// ----

type MaxValueClapmper struct {
	max attr.Value
}

func NewMaxValueClapmper(max attr.Value) *MaxValueClapmper {
	return &MaxValueClapmper{
		max: max,
	}
}

func (c *MaxValueClapmper) Clamp(v attr.Value) attr.Value {
	if v > c.max {
		return c.max
	}
	return v
}

// ----

type MinMaxValueClapmper struct {
	min attr.Value
	max attr.Value
}

func NewMinMaxValueClapmper(min, max attr.Value) *MinMaxValueClapmper {
	return &MinMaxValueClapmper{
		min: min,
		max: max,
	}
}

func (c *MinMaxValueClapmper) Clamp(v attr.Value) attr.Value {
	if v < c.min {
		return c.min
	}
	if v > c.max {
		return c.max
	}
	return v
}
