package firestore

import (
	"concert-manager/data"
	"concert-manager/log"
	"context"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

const artistCollection string = "artists"
var artistFields = []string{"Name", "Genre"}

type ArtistRepo struct {
	Connection *Firestore
}

type ArtistEntity struct {
	Name  string
	Genre string
}

type Artist = data.Artist

func (repo *ArtistRepo) Add(ctx context.Context, artist Artist) (string, error) {
	log.Debug("Attempting to add artist", artist)
	existingArtist, err := repo.findDocRef(ctx, artist.Name)
	if err == nil {
		log.Debugf("Skipping adding artist because it already exists %+v, %v", artist, existingArtist.Ref.ID)
		return existingArtist.Ref.ID, nil
	}
	if err != iterator.Done {
		log.Errorf("Error occurred while checking if artist %v already exists, %v", artist, err)
		return "", err
	}

	artistEntity := ArtistEntity{artist.Name, artist.Genre}
	artists := repo.Connection.Client.Collection(artistCollection)
	docRef, _, err := artists.Add(ctx, artistEntity)
	if err != nil {
		log.Errorf("Failed to add new artist %+v, %v", artist, err)
		return "", err
	}
	log.Infof("Created new artist %+v", docRef.ID)
	return docRef.ID, nil
}

func (repo *ArtistRepo) Update(ctx context.Context, id string, artist Artist) error {
    log.Debug("Attempting to update artist", id, artist)
	docId := repo.Connection.Client.Collection(artistCollection).Doc(id)
	artistDoc, err := docId.Get(ctx)
	if err != nil {
		log.Errorf("Failed to find existing artist while updating %+v, %v", id, err)
		return err
	}

	artistEntity := ArtistEntity{artist.Name, artist.Genre}
	_, err = artistDoc.Ref.Set(ctx, artistEntity)
	if err != nil {
		log.Errorf("Failed to update artist %+v to %v, %v", id, artist, err)
		return err
	}
	log.Info("Successfully updated artist", id)
	return nil
}

func (repo *ArtistRepo) Delete(ctx context.Context, id string) error {
	log.Debug("Attempting to delete artist", id)
	docId := repo.Connection.Client.Collection(artistCollection).Doc(id)
	artistDoc, err := docId.Get(ctx)
	if err != nil {
		log.Errorf("Failed to find existing artist while deleting %+v, %v", id, err)
		return err
	}
	_, err = artistDoc.Ref.Delete(ctx)
	if err != nil {
		log.Error("Failed to delete artist", id, err)
		return err
	}
	log.Infof("Successfully deleted artist %+v", id)
	return nil
}

func (repo *ArtistRepo) Exists(ctx context.Context, artist Artist) (bool, error) {
	log.Debug("Checking for existence of artist", artist)
	doc, err := repo.findDocRef(ctx, artist.Name)
	if err == iterator.Done {
		log.Debug("No existing artist found for", artist)
		return false, nil
	}
	if err != nil {
		log.Errorf("Error while checking existence of artist %v, %v", artist, err)
		return false, err
	}
	log.Debugf("Found artist %v with document ID %v", artist, doc.Ref.ID)
	return true, nil
}

func (repo *ArtistRepo) FindAll(ctx context.Context) ([]Artist, error) {
	log.Debug("Finding all artists")
	artistDocs, err := repo.Connection.Client.Collection(artistCollection).
		Select(artistFields...).
		Documents(ctx).
		GetAll()
	if err != nil {
		log.Error("Error while finding all artists,", err)
		return nil, err
	}

	artists := []Artist{}
	for _, a := range artistDocs {
		artists = append(artists, toArtist(a))
	}
	log.Debugf("Found %d artists", len(artists))
	return artists, nil
}

func toArtist(doc *firestore.DocumentSnapshot) Artist {
    artistData := doc.Data()
	return Artist{
		Name:  artistData["Name"].(string),
		Genre: artistData["Genre"].(string),
		Id:    doc.Ref.ID,
	}
}

func (repo *ArtistRepo) findDocRef(ctx context.Context, name string) (*firestore.DocumentSnapshot, error) {
	artist, err := repo.Connection.Client.Collection(artistCollection).
		Select().
		Where("Name", "==", name).
		Documents(ctx).
		Next()
	if err != nil {
		return nil, err
	}
	return artist, nil
}

func (repo *ArtistRepo) findAllDocs(ctx context.Context) (*map[string]Artist, error) {
	artistDocs, err := repo.Connection.Client.Collection(artistCollection).
		Select(artistFields...).
		Documents(ctx).
		GetAll()
	if err != nil {
		return nil, err
	}

	artists := make(map[string]Artist)
	for _, a := range artistDocs {
		artists[a.Ref.ID] = toArtist(a)
	}

	return &artists, nil
}
