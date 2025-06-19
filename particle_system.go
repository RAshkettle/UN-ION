package main

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Particle struct {
	X, Y       float64
	VX, VY     float64
	Life       float64
	MaxLife    float64
	Size       float64
	R, G, B, A float64
}


type ParticleSystem struct {
	particles []Particle
}


func NewParticleSystem() *ParticleSystem {
	return &ParticleSystem{
		particles: make([]Particle, 0),
	}
}


func (ps *ParticleSystem) applyColorVariation(baseR, baseG, baseB float64) (float64, float64, float64) {
	r := baseR + (rand.Float64()-0.5)*0.3
	g := baseG + (rand.Float64()-0.5)*0.3
	b := baseB + (rand.Float64()-0.5)*0.3

	if r < 0 {
		r = 0
	}
	if r > 1 {
		r = 1
	}
	if g < 0 {
		g = 0
	}
	if g > 1 {
		g = 1
	}
	if b < 0 {
		b = 0
	}
	if b > 1 {
		b = 1
	}

	return r, g, b
}


func (ps *ParticleSystem) AddExplosion(worldX, worldY float64, blockType BlockType) {
	numParticles := 8 + rand.Intn(5) 

	var baseR, baseG, baseB float64
	switch blockType {
	case PositiveBlock:
		baseR, baseG, baseB = 1.0, 0.2, 0.2
	case NegativeBlock:
		baseR, baseG, baseB = 0.2, 0.2, 1.0
	default :
		baseR, baseG, baseB = 0.8, 0.8, 0.8
	}
	for i := 0; i < numParticles; i++ {
		angle := rand.Float64() * 2 * math.Pi

		speed := 20.0 + rand.Float64()*40.0 

		vx := math.Cos(angle) * speed
		vy := math.Sin(angle) * speed

		vy -= 10.0 

		maxLife := 0.3 + rand.Float64()*0.2 
		
		size := 1.5 + rand.Float64()*2.5

		r, g, b := ps.applyColorVariation(baseR, baseG, baseB)

		particle := Particle{
			X:       worldX,
			Y:       worldY,
			VX:      vx,
			VY:      vy,
			Life:    maxLife,
			MaxLife: maxLife,
			Size:    size,
			R:       r,
			G:       g,
			B:       b,
			A:       1.0,
		}

		ps.particles = append(ps.particles, particle)
	}
}

func (ps *ParticleSystem) AddDustCloud(worldX, worldY float64) {
	numParticles := 8 + rand.Intn(6)

	for i := 0; i < numParticles; i++ {
		angle := math.Pi + (rand.Float64()-0.5)*math.Pi*0.8
		speed := 50.0 + rand.Float64()*40.0

		vx := math.Cos(angle) * speed
		vy := math.Sin(angle) * speed

		offsetX := (rand.Float64() - 0.5) * 40.0

		maxLife := 0.8 + rand.Float64()*0.5
		size := 4.0 + rand.Float64()*6.0

		r := 0.9 + rand.Float64()*0.1
		g := 0.9 + rand.Float64()*0.1
		b := 0.9 + rand.Float64()*0.1

		particle := Particle{
			X:       worldX + offsetX,
			Y:       worldY,
			VX:      vx,
			VY:      vy,
			Life:    maxLife,
			MaxLife: maxLife,
			Size:    size,
			R:       r,
			G:       g,
			B:       b,
			A:       0.9,
		}

		ps.particles = append(ps.particles, particle)
	}
}

func (ps *ParticleSystem) Update(dt float64) {
	for i := len(ps.particles) - 1; i >= 0; i-- {
		particle := &ps.particles[i]

		particle.X += particle.VX * dt
		particle.Y += particle.VY * dt

		particle.VY += 200.0 * dt

		resistanceFactor := math.Pow(0.98, dt*60.0)
		particle.VX *= resistanceFactor
		particle.VY *= resistanceFactor

		particle.Life -= dt

		lifeRatio := particle.Life / particle.MaxLife
		particle.A = lifeRatio

		particle.Size = particle.Size * (0.5 + lifeRatio*0.5)

		if particle.Life <= 0 {
			ps.particles[i] = ps.particles[len(ps.particles)-1]
			ps.particles = ps.particles[:len(ps.particles)-1]
		}
	}
}

func (ps *ParticleSystem) Draw(screen *ebiten.Image) {
	for _, particle := range ps.particles {
		if particle.A <= 0 {
			continue
		}

		x := float32(particle.X)
		y := float32(particle.Y)
		radius := float32(particle.Size)

		r := uint8(particle.R * particle.A * 255)
		g := uint8(particle.G * particle.A * 255)
		b := uint8(particle.B * particle.A * 255)
		a := uint8(particle.A * 255)

		vector.DrawFilledCircle(screen, x, y, radius, color.RGBA{r, g, b, a}, true)
	}
}

func (ps *ParticleSystem) HasActiveParticles() bool {
	return len(ps.particles) > 0
}
