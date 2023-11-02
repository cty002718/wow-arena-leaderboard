package controller

import (
	"net/http"
	"sort"
	"time"

	"github.com/cty002718/wow-arena-leaderboard/pkg/ctx"
	"github.com/cty002718/wow-arena-leaderboard/pkg/dao"
	"github.com/gin-gonic/gin"
)

func GetLatestLeaderboard(c *gin.Context) {
	var input struct {
		Season int `form:"season" binding:"required"`
	}

	if err := c.ShouldBindQuery(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}

	leaderboard_dao := ctx.Get[dao.ILeaderboardDao]()
	leaderboard, err := leaderboard_dao.GetLatest(input.Season, "5v5")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}

	arena_record_dao := ctx.Get[dao.IArenaRecordDao]()
	arena_records, err := arena_record_dao.FindByLeaderboardId(leaderboard.ID)

	c.JSON(http.StatusOK, arena_records)
}

type CharacterArenaRecordDto struct {
	CreatedAt   string `json:"created_at"`
	CharacterId int64  `json:"character_id"`
	Season      int    `json:"season"`
	Bracket     string `json:"bracket"`
	Rank        int    `json:"rank"`
	Rating      int    `json:"rating"`
	Won         int    `json:"won"`
	Lost        int    `json:"lost"`
}

func GetCharacterArenaRecord(c *gin.Context) {
	var input struct {
		CharacterId int64  `uri:"character_id" binding:"required"`
		Season      int    `uri:"season" binding:"required"`
		Bracket     string `uri:"bracket" binding:"required"`
	}

	if err := c.ShouldBindUri(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}

	arena_record_dao := ctx.Get[dao.IArenaRecordDao]()
	arena_records, err := arena_record_dao.FindByCharacter(input.CharacterId, input.Season, input.Bracket)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
	}

	sort.Slice(arena_records, func(i, j int) bool {
		return arena_records[i].Leaderboard.CreatedAt.After(arena_records[j].Leaderboard.CreatedAt)
	})

	loc, _ := time.LoadLocation("Asia/Shanghai")
	var arena_records_dto []*CharacterArenaRecordDto
	for _, arena_record := range arena_records {
		arena_records_dto = append(arena_records_dto, &CharacterArenaRecordDto{
			CreatedAt:   arena_record.Leaderboard.CreatedAt.In(loc).Format("2006-01-02 15:04:05"),
			CharacterId: arena_record.CharacterID,
			Season:      arena_record.Leaderboard.SeasonID,
			Bracket:     arena_record.Leaderboard.Bracket,
			Rank:        arena_record.Rank,
			Rating:      arena_record.Rating,
			Won:         arena_record.Won,
			Lost:        arena_record.Lost,
		})
	}

	c.JSON(http.StatusOK, arena_records_dto)
}
