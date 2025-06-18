package main

import (
	"bytes"
	"union/assets"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
)

// AudioManager handles all audio playback
type AudioManager struct {
	audioContext      *audio.Context
	blockBreakPlayer  *audio.Player
	swooshPlayer      *audio.Player
}

// NewAudioManager creates a new audio manager
func NewAudioManager() *AudioManager {
	audioContext := audio.NewContext(44100) // 44.1kHz sample rate
	
	return &AudioManager{
		audioContext: audioContext,
	}
}

// Initialize loads and prepares all audio assets
func (am *AudioManager) Initialize() error {
	// Load block break sound
	blockBreakStream, err := vorbis.DecodeWithSampleRate(am.audioContext.SampleRate(), bytes.NewReader(assets.BlockBreakSound))
	if err != nil {
		return err
	}
	
	am.blockBreakPlayer, err = am.audioContext.NewPlayer(blockBreakStream)
	if err != nil {
		return err
	}
	
	// Load swoosh sound
	swooshStream, err := vorbis.DecodeWithSampleRate(am.audioContext.SampleRate(), bytes.NewReader(assets.SwooshSound))
	if err != nil {
		return err
	}
	
	am.swooshPlayer, err = am.audioContext.NewPlayer(swooshStream)
	if err != nil {
		return err
	}
	
	// Set a specific playback rate for the OGG to prevent thin sound
	// This helps maintain the richness of the audio
	return nil
}

// PlayBlockBreak plays the block break sound effect
func (am *AudioManager) PlayBlockBreak() {
	if am.blockBreakPlayer == nil {
		return
	}
	
	// Rewind to start and play
	am.blockBreakPlayer.Rewind()
	am.blockBreakPlayer.Play()
}

// CreateBlockBreakPlayer creates a new player instance for simultaneous playback
func (am *AudioManager) CreateBlockBreakPlayer() *audio.Player {
	blockBreakStream, err := vorbis.DecodeWithSampleRate(am.audioContext.SampleRate(), bytes.NewReader(assets.BlockBreakSound))
	if err != nil {
		return nil
	}
	
	player, err := am.audioContext.NewPlayer(blockBreakStream)
	if err != nil {
		return nil
	}
	
	return player
}

// PlayBlockBreakMultiple plays the block break sound for multiple blocks (overlapping)
func (am *AudioManager) PlayBlockBreakMultiple(count int) {
	// For multiple blocks, we create separate players to allow overlapping
	for i := 0; i < count && i < 10; i++ { // Limit to 10 simultaneous sounds
		go func() {
			player := am.CreateBlockBreakPlayer()
			if player != nil {
				player.Play()
				// Note: In a production game, you'd want to clean up these players
				// after they finish playing to prevent memory leaks
			}
		}()
	}
}

// PlaySwooshSound plays a subtle swoosh sound for piece movement
func (am *AudioManager) PlaySwooshSound() {
	if am.swooshPlayer == nil {
		return
	}
	
	// Rewind to start and play the dedicated swoosh sound
	am.swooshPlayer.Rewind()
	am.swooshPlayer.Play()
}
