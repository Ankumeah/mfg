# mfg
Simple program in golang to download all your favorite mangas on [mangafreak](https://ww2.mangafreak.me)

## How to download
  ### For regular users
  - Go to [releases](https://github.com/Ankumeah/mfg/releases)
  - Download the version that matches your system (you might have to google some stuff to find out which version)

  ### For chads
  ```sh
  git clone https://github.com/Ankumeah/mfg
  cd mfg
  go build -ldflags="-s -w" -o ./bin/mfg .
  ```

## How to use
- Go to [mangafreak](https://ww2.mangafreak.me/Manga/Chainsaw_Man)
- Navigate to the manga you want to download (*for example - https://ww2.mangafreak.me/Manga/Chainsaw_Man*)
- Copy the last part of the URL (*in this example Chainsaw_Man*)
- Run the command
  - `mfg <your-manga>` to auto discover and download all chapters (in this example `mfg Chainsaw_Man`)
    > Note: Downloaded chapters assuming they have not been renamed or moved, will not be redownloaded (hopefully) so enjoy :)
- A folder with the same name as your manga will appear in your current working directory containing .cbz files

> Note: This project is very much untested so feel free to open any issues if you find any or if you are willing then you can even open a pull request, all (except AI) are welcome :)
