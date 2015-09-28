#Basics
This is a basic little action RPG! Arrow keys to move, space bar to cast your fireball spell, collision with enemies hurts you. It kind of breaks when you die. Also, the ghost of the enemy will still hurt you- oh well! A lot of stuff isn't implemented, but a lot of stuff kind of works.

```bash
go get -u github.com/JoelOtter/termloop
go run game.go
```

[Video of the sprite in action](https://www.youtube.com/watch?v=Dcs2bM05X7I)

#Post Mortem
Termloop does a lot of nice little things right out of the (Term)box. The collision detection works well for a hitbox style. However, it doesn't seem to play well with large background cavases, which was my plan for all of the background art for this game. It gets pretty slow pretty quickly (if the background forest artwork is removed, the game plays a lot more smoothly). Overall, I think I tried to make too complicated of a game to really finish in just a weekend- but I learned a lot about Go and had some fun!

#Credits
Character and enemy sprite artwork is by Svetlana Kushnariova (lana-chan@yandex.ru): http://opengameart.org/content/24x32-characters-with-faces-big-pack

Forest tiles artwork by Stephen Challener: http://opengameart.org/content/32x32-and-16x16-rpg-tiles-forest-and-some-interior-tiles

Fireball artwork is by [Clint Bellanger](http://clintbellanger.net): http://opengameart.org/content/fireball-spell
