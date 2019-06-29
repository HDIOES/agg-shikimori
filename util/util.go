package util

type Configuration struct {
	DatabaseURL             string `json:"databaseUrl"`
	MaxOpenConnections      int    `json:"maxOpenConnections"`
	MaxIdleConnections      int    `json:"maxIdleConnections"`
	ConnectionTimeout       int    `json:"connectionTimeout"`
	Port                    int    `json:"port"`
	ShikimoriURL            string `json:"shikimori_url"`
	ShikimoriAnimeSearchURL string `json:"shikimori_anime_search_url"`
	ShikimoriGenreURL       string `json:"shikimori_genre_url"`
	ShikimoriStudioURL      string `json:"shikimori_studio_url"`
}
