
# PhotoBook Server (pbs)

pbs is a small project demonstrating how:

 - [github.com/clarktrimble/delish](https://github.com/clarktrimble/delish)
 - [github.com/clarktrimble/giant](https://github.com/clarktrimble/giant)
 - [github.com/clarktrimble/launch](https://github.com/clarktrimble/launch)
 - [github.com/clarktrimble/sabot](https://github.com/clarktrimble/sabot)

can work together to add industrial strength to a Golang project.

## pbs Components

pbs consists of:

### Resize

A utility to resize photos from a Google Takeout tarball.

### Load

A utility to post metadata from a Takeout tarball and resized images to pbs api.

### Api

A json web api supporting a PhotoBook App.

## The Missing App

The PhotoBook App is yet to be published.  Suffice it to say that it consumes objects of the form:

```json
{
  "photo_id": "7hwDrAx",
  "src": "http://tartu/photo/resized/PXL_20230722_090048284-large.png",
  "width": 1020,
  "height": 768,
  "thumb": "http://tartu/photo/resized/PXL_20230722_090048284-thumb.png",
  "thumb_gs": "http://tartu/photo/resized/PXL_20230722_090048284-thumb-gs.png",
  "lat": 60.3149083,
  "lon": 24.9611861,
  "taken_at": "2023-07-22T09:00:48Z",
  "featured": true
}
```

## Salient Points

Salacious details forthcoming!

In the meantime, cmd/api provides a good starting point for the curious:

```go
type Config struct {
  Version  string         `json:"version" ignored:"true"`
  Truncate int            `json:"truncate" desc:"truncate log fields beyond length"`
  Bolt     *bolt.Config   `json:"bolt"`
  Server   *delish.Config `json:"server"`
}

func main() {

  // load config, setup logger

  cfg := &Config{Version: version}
  launch.Load(cfg, cfgPrefix)

  lgr := &sabot.Sabot{Writer: os.Stdout, MaxLen: cfg.Truncate}
  ctx := lgr.WithFields(context.Background(), "run_id", hondo.Rand(7))

  ctx = graceful.Initialize(ctx, &wg, lgr)

  // create router/handler, and server

  rtr := chi.New()

  handler := mid.LogResponse(lgr, rtr)
  handler = mid.LogRequest(lgr, hondo.Rand, handler)
  handler = mid.ReplaceCtx(ctx, handler)

  svr := cfg.Server.New(handler, lgr)

  // setup service layer and register routes

  repo, err := cfg.Bolt.New()
  launch.Check(ctx, lgr, err)
  defer repo.Close()

  photoSvc := &photosvc.PhotoSvc{
          Logger: lgr,
          Repo:   repo,
  }

  photoSvc.Register(rtr)
  rtr.Set("GET", "/config", svr.ObjHandler("config", cfg))

  // delicious!

  svr.Start(ctx, &wg)
  graceful.Wait(ctx)
}
```

## Disclaimer

The above demonstrates how one could choose to follow a well-trodden path to implement http api's and clients in Golang.
Factoring and reusing some parts of a well-trodden path can, eventually, save time and increase reliability.
Small, lightweight, and well constrained modules tend to break-even more often and more quickly.

The underlying approaches illustrated have some general acceptance.
B-but, neither the well-trodden path presented nor the underlying approches illustrated are considered superior to others.

One could, ever so gently, make the case that a team choosing and cohering to some set workable paths/approaches is likely to outperform those that do not.
In the same breath, one could also reflect on the inefficiencies of large-scale and/or orthodox proscriptions we may have had the misfortune to toil beneath.

So yeah, it's about code reuse at the tactical level, which can be fun to think about and perhaps even profitable :)

## Alienated Todo's

```text
    photosvc/photosvc.go
    16:// Todo: all photos are global, need a scoping concept, i.e.: "Baltic Travels"

    resize/resize.go
       // Todo: comment on name to filename conventions somewhere plz (thumb, large, and so forth)
    17:// Todo: short post on mutating Golang slices, may need to look back in log..
    18:// Todo: look at pulling straight from  takeout from tarball
    19:// Todo: look at using multiple cores
    65:// Todo: would save time to combo yeah? (gs that is)

    other
    Always nice to log a copy of the loaded config right away!
    Build test instructions
```
