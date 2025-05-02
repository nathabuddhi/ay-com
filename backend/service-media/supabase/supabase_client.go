package supabase

import (
	"os"

	storage_go "github.com/supabase-community/storage-go"
)

var Client *storage_go.Client

func InitSupabase() {
	url := os.Getenv("SUPABASE_URL")
	key := os.Getenv("SUPABASE_KEY")

	storageURL := url + "/storage/v1"

	Client = storage_go.NewClient(storageURL, key, nil)
}
