package doarama_test

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/twpayne/go-doarama"
)

func Example() {
	// Create the client using anonymous authentication
	client := doarama.NewClient(
		doarama.APIKey("xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"),
		doarama.APIName("Your API Name"),
		doarama.Anonymous("userid"),
	)

	// Open the GPS track
	filename := "activity GPS filename (GPX or IGC)"
	gpsTrack, err := os.Open(filename)
	if err != nil {
		return
	}
	defer gpsTrack.Close()

	ctx := context.Background()

	// Create the activity
	activity, err := client.CreateActivity(ctx, filepath.Base(filename), gpsTrack)
	if err != nil {
		return
	}
	log.Printf("ActivityId: %d", activity.ID)

	// Set the activity info
	if err := activity.SetInfo(ctx, &doarama.ActivityInfo{
		TypeID: doarama.FlyParaglide,
	}); err != nil {
		return
	}

	// Create the visualisation
	activities := []*doarama.Activity{activity}
	visualisation, err := client.CreateVisualisation(ctx, activities)
	if err != nil {
		return
	}
	log.Printf("VisualisationKey: %s", visualisation.Key)
	log.Printf("VisualisationURL: %s", visualisation.URL(nil))
}
