package main

import "github.com/hajimehoshi/ebiten/v2"

type SceneType int

const (
	SceneTitleScreen SceneType = iota
	SceneGame
	SceneEndScreen
	SceneHelp // Add help scene type
)

type Scene interface {
	Update() error
	Draw(screen *ebiten.Image)
	Layout(outerWidth, outerHeight int) (int, int)
}

type SceneManager struct {
	currentScene Scene
	sceneType    SceneType
	titleScene   *TitleScene
	gameScene    *GameScene
	endScene     *EndScene
	helpScene    *HelpScene
}

func (sm *SceneManager) Update() error {
	return sm.currentScene.Update()
}

func (sm *SceneManager) Draw(screen *ebiten.Image) {
	sm.currentScene.Draw(screen)
}

func (sm *SceneManager) Layout(outerWidth, outerHeight int) (int, int) {
	return sm.currentScene.Layout(outerWidth, outerHeight)
}
func NewSceneManager() *SceneManager {
	sm := &SceneManager{
		sceneType: SceneTitleScreen,
	}

	sm.titleScene = NewTitleScene(sm)
	sm.gameScene = NewGameScene(sm)
	sm.endScene = NewEndScene(sm, 0)
	sm.helpScene = NewHelpScene(sm)

	sm.currentScene = sm.titleScene

	return sm
}

func (sm *SceneManager) TransitionTo(sceneType SceneType) {
	sm.sceneType = sceneType

	switch sceneType {
	case SceneTitleScreen:
		sm.currentScene = sm.titleScene
		sm.titleScene.prevHPressed = true
	case SceneGame:
		sm.currentScene = sm.gameScene
	case SceneEndScreen:
		sm.currentScene = sm.endScene
	case SceneHelp:
		sm.currentScene = sm.helpScene
		sm.helpScene.prevHPressed = true
	}
}

func (sm *SceneManager) TransitionToEndScreen(finalScore int) {
	sm.sceneType = SceneEndScreen
	sm.endScene = NewEndScene(sm, finalScore)
	sm.currentScene = sm.endScene
}

func (sm *SceneManager) GetCurrentSceneType() SceneType {
	return sm.sceneType
}
