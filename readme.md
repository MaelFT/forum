COMBY Quentin
DIEDERICHS Geoffrey
FATH Mael

# Projet Forum YNOV informatique B1

## Utilisation

- Depuis la racine du projet lancer le server en saisissant dans votre terminal :

```sh
    $ go run main.go
```

- Lancez votre navigateur web et allez sur la page localhost:8080

- Vous pouvez accéder à certaines pages du projet telles-que :
    /index
    /404
    /login
    /register
    /category
    /post?id=1 ou /post?id=2 ou /post?id=3

## Organisation des fichiers

- ~/controllers contient le back du projet : la base de donnée, et le code gérant les pages webs

- ~/models contient les structures json utilisées pour transmettre des données

- ~/views contient le front : les pages html, le css (~/views/assets/style.css), les images utilisées dans le front (~/views/assets/images/) et les polices (~/views/assets/fonts/)

- ~/main.go lance le serveur web