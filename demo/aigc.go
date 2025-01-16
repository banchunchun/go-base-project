package demo

type MediaPlayer interface {
	Play()
}

type VideoPlayer struct {
	VideoName string
}

func (p *VideoPlayer) Play() {
	p.Play()
}
