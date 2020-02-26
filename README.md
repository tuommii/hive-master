# Hive Master
<img src="https://github.com/wehard/hive-master/blob/master/assets/hive_master.png?raw=true"/>

## Build

### Install SDL2

#### Linux
```sudo apt-get install libsdl2{,-image,-mixer,-ttf,-gfx}-dev```

#### OS X
```brew install sdl2{,_image,_mixer,_ttf,_gfx} pkg-config```

### Install Go Stuff

```
go get -v github.com/veandco/go-sdl2/sdl
go get -v github.com/veandco/go-sdl2/img
go get -v github.com/veandco/go-sdl2/mix
go get -v github.com/veandco/go-sdl2/ttf
go get -v github.com/veandco/go-sdl2/gfx

go get -d "github.com/wehard/hive-master/"
```

### Run
```cd $GOPATH/src/github.com/wehard/hive-master```

```go build```

```./hive-master```
