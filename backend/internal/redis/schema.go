package redis

const (
	// Cache keys
	MangaDetailKey    = "manga:detail:%s"     // manga:detail:<manga_id>
	ChapterDetailKey  = "chapter:detail:%s"    // chapter:detail:<chapter_id>
	ChapterContentKey = "chapter:content:%s"   // chapter:content:<chapter_id>
	
	// Counter keys
	MangaViewKey     = "manga:views:%s"       // manga:views:<manga_id>
	ChapterViewKey   = "chapter:views:%s"      // chapter:views:<chapter_id>
	
	// Sorted sets
	PopularMangaKey  = "manga:popular"         // Sorted set of manga by view count
	LatestMangaKey   = "manga:latest"          // Sorted set of manga by update time
	
	// Cache expiration times
	MangaDetailTTL   = 3600  // 1 hour
	ChapterDetailTTL = 3600  // 1 hour
	ContentCacheTTL  = 7200  // 2 hours
)

// ViewCountUpdate increments view count and updates sorted sets
func ViewCountUpdate(rdb *redis.Client, mangaID, chapterID string) error {
	pipe := rdb.Pipeline()
	
	// Increment manga views
	pipe.Incr(fmt.Sprintf(MangaViewKey, mangaID))
	
	// Increment chapter views if provided
	if chapterID != "" {
		pipe.Incr(fmt.Sprintf(ChapterViewKey, chapterID))
	}
	
	// Update manga score in popular list
	pipe.ZIncrBy(PopularMangaKey, 1, mangaID)
	
	_, err := pipe.Exec()
	return err
}
