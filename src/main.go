package main

import (
	"flag"
	"fmt"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/auth"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"image"
	"image/png"
	"os"
)

var (
	ip   string
	port string
)

// ty TwistedAsylum in the gophertunnel discord
func SkinToRGBA(s protocol.Skin) *image.RGBA {
	t := image.NewRGBA(image.Rect(0, 0, int(s.SkinImageWidth), int(s.SkinImageHeight)))
	t.Pix = s.SkinData
	return t
}

func antiCrash(p protocol.PlayerListEntry) bool {
	if p.Skin.SkinData == nil {
		return true
	} else {
		return false
	}
}

func main() {
	flag.StringVar(&ip, "ip", "127.0.0.1", "Servers IP Address")
	flag.StringVar(&port, "port", "19132", "Servers Port")
	flag.Parse()
	_ = os.Mkdir("stolen", 0755)

	dialer := minecraft.Dialer {
		TokenSource: auth.TokenSource,
	}

	address := ip + ":" + port
	fmt.Println("Connecting to " + address)
	conn, err := dialer.Dial("raknet", address)
	if err != nil {
		panic(err)
	}

	defer conn.Close()
	_ = conn.DoSpawn()

	for {
		pk, err := conn.ReadPacket()
		if err != nil {
			break
		}

		switch p := pk.(type) {
		case *packet.PlayerList:
			go func() {
				for _, player := range p.Entries {
					if antiCrash(player) {
						return
					}
					name := player.Username
					skin := SkinToRGBA(player.Skin)
					path, _ := os.Getwd()
					fileSkin, _ := os.Create(fmt.Sprintf("%s/stolen/%s.png", path, name))
					_ = png.Encode(fileSkin, skin)
					fileSkin.Close()
					fmt.Println("Stolen " + name)
				}
			}()
		}
	}
}
