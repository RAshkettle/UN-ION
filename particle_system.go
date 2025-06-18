package main

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Particle represents a single particle in an explosion
type Particle struct {
	X, Y       float64 // Current position
	VX, VY     float64 // Velocity
	Life       float64 // Life remaining (0-1)
	MaxLife    float64 // Initial life
	Size       float64 // Particle size
	R, G, B, A float64 // Color components
}

// ParticleSystem manages all particles and explosions
type ParticleSystem struct {
	particles []Particle
}

// NewParticleSystem creates a new particle system
func NewParticleSystem() *ParticleSystem {
	return &ParticleSystem{
		particles: make([]Particle, 0),
	}
}

// AddExplosion creates a cartoony explosion at the specified world coordinates
func (ps *ParticleSystem) AddExplosion(worldX, worldY float64, blockType BlockType) {
	numParticles := 8 + rand.Intn(5) // 8-12 particles (reduced from 15-25)
	
	// Get base color based on block type
	var baseR, baseG, baseB float64
	switch blockType {
	case PositiveBlock:
		baseR, baseG, baseB = 1.0, 0.2, 0.2 // Red-ish
	case NegativeBlock:
		baseR, baseG, baseB = 0.2, 0.2, 1.0 // Blue-ish
	case NeutralBlock:
		baseR, baseG, baseB = 0.8, 0.8, 0.8 // Gray-ish
	default:
		baseR, baseG, baseB = 1.0, 1.0, 1.0 // White
	}
	
	for i := 0; i < numParticles; i++ {
		// Random angle for explosion direction
		angle := rand.Float64() * 2 * math.Pi
		
		// Reduced speed to keep particles more contained to the block
		speed := 20.0 + rand.Float64()*40.0 // Reduced from 50-150 to 20-60
		
		// Calculate velocity components
		vx := math.Cos(angle) * speed
		vy := math.Sin(angle) * speed
		
		// Reduced upward bias
		vy -= 10.0 // Reduced from 20.0
		
		// Shorter life duration to keep explosion more contained
		maxLife := 0.3 + rand.Float64()*0.2 // Reduced from 0.4-0.6 to 0.3-0.5
		
		// Smaller particle size
		size := 1.5 + rand.Float64()*2.5 // Reduced from 2-6 to 1.5-4
		
		// Color variation
		r := baseR + (rand.Float64()-0.5)*0.3
		g := baseG + (rand.Float64()-0.5)*0.3
		b := baseB + (rand.Float64()-0.5)*0.3
		
		// Clamp colors
		if r < 0 { r = 0 }
		if r > 1 { r = 1 }
		if g < 0 { g = 0 }
		if g > 1 { g = 1 }
		if b < 0 { b = 0 }
		if b > 1 { b = 1 }
		
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

// AddDustCloud creates a small dust cloud when a piece hits the bottom
func (ps *ParticleSystem) AddDustCloud(worldX, worldY float64) {
	numParticles := 8 + rand.Intn(6) // 8-13 particles (more visible)
	
	for i := 0; i < numParticles; i++ {
		// Horizontal spread - wider spread for more visibility
		angle := math.Pi + (rand.Float64()-0.5)*math.Pi*0.8 // Wider spread
		speed := 50.0 + rand.Float64()*40.0 // Faster for more visibility
		
		vx := math.Cos(angle) * speed
		vy := math.Sin(angle) * speed
		
		// Wider random offset from center
		offsetX := (rand.Float64() - 0.5) * 40.0
		
		// Dust properties - more pronounced
		maxLife := 0.8 + rand.Float64()*0.5 // Longer lived
		size := 4.0 + rand.Float64()*6.0 // Bigger particles
		
		// White dust color for natural look
		r := 0.9 + rand.Float64()*0.1 // Light gray/white
		g := 0.9 + rand.Float64()*0.1 // Light gray/white
		b := 0.9 + rand.Float64()*0.1 // Light gray/white
		
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
			A:       0.9, // More opaque for visibility
		}
		
		ps.particles = append(ps.particles, particle)
	}
}

// Update updates all particles
func (ps *ParticleSystem) Update(dt float64) {
	// Update existing particles
	for i := len(ps.particles) - 1; i >= 0; i-- {
		particle := &ps.particles[i]
		
		// Update position
		particle.X += particle.VX * dt
		particle.Y += particle.VY * dt
		
		// Apply gravity
		particle.VY += 200.0 * dt
		
		// Apply air resistance
		particle.VX *= 0.98
		particle.VY *= 0.98
		
		// Update life
		particle.Life -= dt
		
		// Update alpha based on remaining life
		lifeRatio := particle.Life / particle.MaxLife
		particle.A = lifeRatio
		
		// Update size (shrink over time)
		particle.Size = particle.Size * (0.5 + lifeRatio*0.5)
		
		// Remove dead particles
		if particle.Life <= 0 {
			// Remove particle by swapping with last and shrinking slice
			ps.particles[i] = ps.particles[len(ps.particles)-1]
			ps.particles = ps.particles[:len(ps.particles)-1]
		}
	}
}

// Draw renders all particles
func (ps *ParticleSystem) Draw(screen *ebiten.Image) {
	for _, particle := range ps.particles {
		if particle.A <= 0 {
			continue
		}
		
		// Create a simple circular particle using vector drawing
		x := float32(particle.X)
		y := float32(particle.Y)
		radius := float32(particle.Size)
		
		// Convert color to 0-255 range
		r := uint8(particle.R * particle.A * 255)
		g := uint8(particle.G * particle.A * 255)
		b := uint8(particle.B * particle.A * 255)
		a := uint8(particle.A * 255)
		
		// Draw a filled circle for the particle
		vector.DrawFilledCircle(screen, x, y, radius, color.RGBA{r, g, b, a}, true)
	}
}

// HasActiveParticles returns true if there are any active particles
func (ps *ParticleSystem) HasActiveParticles() bool {
	return len(ps.particles) > 0
}
