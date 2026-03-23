# TRMNL-Farside

TRMNL-Farside is a small Go service that scrapes the daily comic from https://www.thefarside.com/ and returns JSON for a TRMNL Private Plugin.

![TRMNL Far Side](FarSide.png)

## What this service returns

The endpoint returns JSON shaped like:

```json
{
  "img": "https://featureassets.amuniversal.com/assets/...",
  "imageUrl": "https://featureassets.amuniversal.com/assets/...",
  "title": "Comic caption or alt text"
}
```

Use `img` in your TRMNL markup.

## Run your own scraping server (optional)

1. Install Go: https://golang.org/doc/install
2. Clone this repository and move into the project folder.
3. Initialize modules (first time only):
   ```bash
   go mod init farside_go
   ```
4. Install dependencies:
   ```bash
   go get github.com/PuerkitoBio/goquery
   go mod tidy
   ```
5. Run the server:
   ```bash
   go run main.go
   ```
6. Expose the server publicly on port `8123`.

## Create the TRMNL Private Plugin

1. Log in to TRMNL and open Plugins.
2. Open Private Plugin and click Add new.
3. Name your plugin (for example, Daily Comic).
4. In Strategy, choose Polling.
5. Set Polling URL to your server URL.
6. Save once so Edit Markup becomes available.

## TRMNL Markup

Paste this into Edit Markup:

```html
<div class="view comic-view bg-white">
  <div class="layout layout--col gap--space-between">
    <div class="columns">
      <div class="column">
        <div class="content content" style="max-height: 95hv">
          <img src="{{ img }}" alt="{{ title }}" style="display: block; margin: 0 auto; max-width: 80%; height: 350px; border-radius: 5px;" />
          <p style="text-align:center">{{ title }}</p>
        </div>
      </div>
    </div>
  </div>
  <div class="title_bar">
    <img class="image" src="https://siteassets.thefarside.com/packs/media/images/brand/logo/tfs_logo-cb7a52999a9300e203e36865dc428f71.svg" alt="The Far Side by Gary Larson" />
  </div>
</div>
```

## Activate and view

1. Save the plugin.
2. Open Playlists in TRMNL.
3. Move the plugin to the desired position.
4. Refresh your TRMNL device to view the daily comic.

## License

This project is licensed under GPL-3.0. See [LICENSE](LICENSE).
