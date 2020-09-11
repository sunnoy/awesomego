package main

import "fmt"

type song struct {
	name     string
	artist   string
	previous *song
	next     *song
}

type playlist struct {
	name       string
	head       *song
	tail       *song
	nowPlaying *song
}

func createPlaylist(name string) *playlist {
	return &playlist{
		name: name,
	}
}

func (p *playlist) addSong(name, ariist string) error {
	s := &song{
		name:   name,
		artist: ariist,
	}

	if p.head == nil {
		p.head = s
	} else {
		//确定当前的最后一个node
		currentNode := p.tail
		currentNode.next = s

		// 补齐song的结构体
		s.previous = p.tail
	}

	// 最新的列表的最后一个
	p.tail = s

	return nil
}

func (p *playlist) showAllSongs() error {
	currentNode := p.head
	if currentNode == nil {
		fmt.Println("playlist is nil")
		return nil
	}
	// 第一首歌曲
	fmt.Printf("%v\n", *currentNode)

	// 下面的几首歌曲进行展示
	for currentNode.next != nil {
		currentNode = currentNode.next
		fmt.Printf("%v\n", *currentNode)
	}
	return nil

}

func (p *playlist) startPlaying() *song {
	p.nowPlaying = p.head
	return p.nowPlaying
}

func (p *playlist) nextSong() *song {
	p.nowPlaying = p.nowPlaying.next
	return p.nowPlaying
}

func (p *playlist) previousSong() *song {
	p.nowPlaying = p.nowPlaying.previous
	return p.nowPlaying
}

func main() {
	playlistName := "myplaylist"

	myPlaylist := createPlaylist(playlistName)
	fmt.Println("Created playlist")
	fmt.Println()

	fmt.Print("Adding songs to the playlist...\n\n")
	myPlaylist.addSong("1", "The Lumineers")
	myPlaylist.addSong("2", "Ed Sheeran")
	myPlaylist.addSong("3", "The Lumineers")
	myPlaylist.addSong("4", "Calvin Harris")

	fmt.Println("Showing all songs in playlist...")
	myPlaylist.showAllSongs()
	fmt.Println()

	fmt.Println("startPlaying")
	myPlaylist.startPlaying()
	fmt.Printf("now playing name %v,art %v", myPlaylist.nowPlaying.name, myPlaylist.nowPlaying.artist)

	fmt.Println()

	fmt.Println("netxPlaying")
	myPlaylist.nextSong()
	fmt.Printf("now playing name %v,art %v", myPlaylist.nowPlaying.name, myPlaylist.nowPlaying.artist)

	fmt.Println()
	fmt.Println("prePlaying")
	myPlaylist.previousSong()
	fmt.Printf("now playing name %v,art %v", myPlaylist.nowPlaying.name, myPlaylist.nowPlaying.artist)

}
