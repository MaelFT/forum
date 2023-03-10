package forum

import (
	"database/sql"
	"fmt"
	models "forum/models"
)

func Like(Like models.Like, likevalue int) {
	db, err := sql.Open("sqlite3", "forum.db")
	if err != nil {
		fmt.Println(err)
	}
	forumRepository := NewSQLiteRepository(db)
	hasLiked, err := userHasLikedPost(db, int(Like.User_ID), int(Like.Post_ID))
	if err != nil {
		panic(err)
	}
	hasDisLiked, erro := userHasDisLikedPost(db, int(Like.User_ID), int(Like.Post_ID))

	if hasDisLiked {
		fmt.Println("l'utilisateur à déjà like le post")
		if likevalue == 1 {
			err = deleteLikesByUser(db, int(Like.User_ID))
			if err != nil {
				panic(erro)
			}
			forumRepository.CreateLike(Like)
		}
		if likevalue == -1 {
			err = deleteLikesByUser(db, int(Like.User_ID))
			if err != nil {
				panic(erro)
			}
		}
	} else if hasLiked {
		fmt.Println("l'utilisateur à déjà like le post")
		if likevalue == 1 {
			err = deleteLikesByUser(db, int(Like.User_ID))
			if err != nil {
				panic(err)
			}
		}
		if likevalue == -1 {
			err = deleteLikesByUser(db, int(Like.User_ID))
			if err != nil {
				panic(err)
			}
			forumRepository.CreateLike(Like)
		}
	} else {
		forumRepository.CreateLike(Like)
	}
}
