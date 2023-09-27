package tui

import (
    "github.com/thomasjinlo/gochatter/internal/client"
)

type Renderer interface {
    Render(client *client.Client)
}

type RendererFunc func(client *client.Client)

func (f RendererFunc) Render(client *client.Client) {
    f(client)
}
