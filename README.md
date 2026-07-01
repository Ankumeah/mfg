# mfg
Simple program in golang to download all your favorite mangas on [mangafreak](https://ww2.mangafreak.me)

> [!WARNING]
> I did not think much of windows users while makeing this so Im sorry in advance if you happen to be on windows

## Build
  ```sh
  git clone https://github.com/Ankumeah/mfg
  cd mfg
  go build -ldflags="-s -w" -o ./bin/mfg main.go
  ```

## How to use
- Go to [mangafreak](https://ww2.mangafreak.me/)
- Navigate to the manga you want to download (*for example - https://ww2.mangafreak.me/Manga/Chainsaw_Man*)
- Copy the last part of the URL (*in this example Chainsaw_Man*)
- Run the command
  - `mfg <your-manga>` to auto discover and download all chapters (in this example `mfg Chainsaw_Man`)
    > Note: Downloaded chapters assuming they have not been renamed or moved, will not be redownloaded (hopefully) so enjoy :)
- A folder with the same name as your manga will appear in your current working directory containing .cbz files

## Flags

|    Flag    |  Type  | Default |                    Discription                    |
|------------|--------|---------|---------------------------------------------------|
| -j         |  uint  |   10    | Max parallel downloads                            |
| -t         | string |   10s   | Max timeout. Eg - 1h1m1s1ms1us1ns (default "10s") |
| -h, --help |        |         | Prints the help text                              |


## helper.sh
This repo contains a helper.sh script which will help manage your library but keep in mind
that this is a script that I quickly hacked together

  ### How to use

  `./helper.sh [mode] <input>`


  | Mode |             Discription             |
  |------|-------------------------------------|
  | a    | Add a manga to the save             |
  | rm   | Remove a manga from the save        |
  | ls   | List mangas from the save           |
  | ia   | Add a file to the ignore            |
  | irm  | Remove a file from the ignore       |
  | ils  | List mangas from the ignore         |
  | s    | Sync to the save                    |
  | c    | Remove files not in save and ignore |

> [!NOTE]
> This project is very much untested so feel free to open any issues if you
find any or if you are willing then you can even open a pull request, all (except AI) are welcome :)
