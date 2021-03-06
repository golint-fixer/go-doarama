# doarama.go

doarama is a command line interface to the [Doarama GPS visualiszation
platform](http://www.doarama.com/).  Install it with:

    $ go get github.com/twpayne/go-doarama/cmd/dorama

## How to set your Doarama API credentials and user id

You must set your Doarama API credentials before running the `doarama` command.

    $ export DOARAMA_API_KEY="Your Doarama API key"
    $ export DOARAMA_API_NAME="Your Doarama API name"
    $ export DOARAMA_USER_ID="Your Doarama user id"

## How to create a visualisation URL of one or more activities in a single step

    $ doarama create \
          --activitytype=paraglide \
          --name="Tom Payne" \
          --name="Christian Erne" \
          2015-04-12-XCT-TPA-01.igc \
          2015-04-12-XCT-CER-01.igc
    ActivityId: 479145
    ActivityId: 479146
    VisualisationKey: E2PKx1e
    VisualisationURL: https://api.doarama.com/api/0.2/visualisation?k=E2PKx1e&name=Tom+Payne&name=Christian+Erne

## How to create a visualisation URL of a single activity step by step

Upload an activity and set its [activity type
id](https://api.doarama.com/api/0.2/activityType):

    $ doarama activity create --activitytype paraglide 2015-08-02-FLY-5094-01.IGC
    ActivityId: 479049

You can specify multiple tracklog files on the command line to create multiple
activities simultaneously.

Create a visualisation of this activity:

    $ doarama visualisation create 479049
    VisualisationKey: eBB1Gwe

You can specify multiple activity ids to create a visualisation with multiple
activities.

Create a URL for this visualisation:

    $ doarama visualisation url --name="Tom Payne" eBB1Gwe
    VisualisationURL: https://api.doarama.com/api/0.2/visualisation?k=eBB1Gwe&name=Tom+Payne
