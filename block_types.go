package main

type BlockType int

const (
	PositiveBlock BlockType = iota
	NegativeBlock
	NeutralBlock
)

type Block struct {
	X, Y          int
	BlockType     BlockType
	IsWobbling    bool
	WobbleTime    float64
	WobblePhase   float64
	ShowPowSprite bool
	IsInStorm     bool
	StormTime     float64
	StormPhase    float64
	SparkPhase    float64
	IsFalling     bool
	FallStartY    float64
	FallTargetY   float64
	FallProgress  float64
	IsArcing      bool
	ArcStartX     float64
	ArcStartY     float64
	ArcTargetX    float64
	ArcTargetY    float64
	ArcProgress   float64
	ArcRotation   float64
	ArcScale      float64
}
