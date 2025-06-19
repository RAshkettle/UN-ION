package main

import (
	"bytes"
	"time"
	"union/assets"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
)

const (
	SampleRate            = 44100
	BackgroundMusicVolume = 0.1
	SoundEffectVolume     = 1.0
)

type AudioManager struct {
	audioContext          *audio.Context
	blockBreakPlayer      *audio.Player
	swooshPlayer          *audio.Player
	backgroundMusicPlayer *audio.Player
	musicPaused           bool
}

func NewAudioManager() *AudioManager {
	audioContext := audio.NewContext(SampleRate)
	

	return &AudioManager{
		audioContext: audioContext,
	}
}

func (am *AudioManager) Initialize() error {
	blockBreakStream, err := vorbis.DecodeWithSampleRate(am.audioContext.SampleRate(), bytes.NewReader(assets.BlockBreakSound))
	if err != nil {
		return err
	}

	am.blockBreakPlayer, err = am.audioContext.NewPlayer(blockBreakStream)
	if err != nil {
		return err
	}
	am.blockBreakPlayer.SetBufferSize(SampleRate / 10)

	swooshStream, err := vorbis.DecodeWithSampleRate(am.audioContext.SampleRate(), bytes.NewReader(assets.SwooshSound))
	if err != nil {
		return err
	}

	am.swooshPlayer, err = am.audioContext.NewPlayer(swooshStream)
	if err != nil {
		return err
	}
	am.swooshPlayer.SetBufferSize(100)

	backgroundMusicStream, err := mp3.DecodeWithSampleRate(am.audioContext.SampleRate(), bytes.NewReader(assets.BackgroundMusic))
	if err != nil {
		return err
	}

	am.backgroundMusicPlayer, err = am.audioContext.NewPlayer(backgroundMusicStream)
	if err != nil {
		return err
	}
	am.backgroundMusicPlayer.SetBufferSize(100)

	am.backgroundMusicPlayer.SetVolume(BackgroundMusicVolume)

	return nil
}

func (am *AudioManager) PlayBlockBreak() {
	if am.blockBreakPlayer == nil {
		return
	}

	am.blockBreakPlayer.Rewind()
	am.blockBreakPlayer.Play()
}

func (am *AudioManager) CreateBlockBreakPlayer() *audio.Player {
	blockBreakStream, err := vorbis.DecodeWithSampleRate(am.audioContext.SampleRate(), bytes.NewReader(assets.BlockBreakSound))
	if err != nil {
		return nil
	}

	player, err := am.audioContext.NewPlayer(blockBreakStream)
	if err != nil {
		return nil
	}
	player.SetBufferSize(100)

	return player
}

func (am *AudioManager) PlayBlockBreakMultiple(count int) {
	for i := 0; i < count && i < 10; i++ {
		go func() {
			player := am.CreateBlockBreakPlayer()
			if player != nil {
				player.Play()
			}
		}()
	}
}

func (am *AudioManager) PlaySwooshSound() {
	if am.swooshPlayer == nil {
		return
	}

	am.swooshPlayer.Rewind()
	am.swooshPlayer.Play()
}

func (am *AudioManager) StartBackgroundMusic() {
	if am.backgroundMusicPlayer == nil {
		return
	}

	am.musicPaused = false
	am.backgroundMusicPlayer.Play()

	go func() {
		for {
			if !am.backgroundMusicPlayer.IsPlaying() && !am.musicPaused {
				am.backgroundMusicPlayer.Rewind()
				am.backgroundMusicPlayer.Play()
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()
}

func (am *AudioManager) StopBackgroundMusic() {
	if am.backgroundMusicPlayer != nil {
		am.musicPaused = true
		am.backgroundMusicPlayer.Pause()
	}
}

func (am *AudioManager) PauseBackgroundMusic() {
	if am.backgroundMusicPlayer != nil {
		am.musicPaused = true
		am.backgroundMusicPlayer.Pause()
	}
}

func (am *AudioManager) ResumeBackgroundMusic() {
	if am.backgroundMusicPlayer != nil {
		am.musicPaused = false
		am.backgroundMusicPlayer.Play()
	}
}
