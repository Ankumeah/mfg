#!/bin/sh

save_file=".save"
ignore=".ignore"
bin="./.mfg"

add() {
  echo "$1" >> "${save_file}"
}
ignore_add() {
  echo "$1" >> "${ignore}"
}
remove() {
  sed -i "/^$1\$/d" "${save_file}"
}
ignore_remove() {
  sed -i "/^$1\$/d" "${ignore}"
}
list_save() {
  cat "${save_file}"
}
list_ignore() {
  cat "${ignore}"
}
sync() {
  cat "${save_file}" | xargs "${bin}"
}
clean() {
  for file in $(ls); do
    grep -qw "${file}" "${ignore}" && continue
    if ! grep -q "^${file}" "${save_file}"; then
      rm -rf "${file}"
    fi
  done
}
help() {
  echo 'a   : Add a manga to the save'
  echo 'rm  : Remove a manga from the save'
  echo 'ls  : List mangas from the save'
  echo 'ia  : Add a file to the ignore'
  echo 'irm : Remove a file from the ignore'
  echo 'ils : List mangas from the ignore'
  echo 's   : Sync to the save'
  echo 'c   : Remove files not in save and ignore'
}

case "$1" in
  "a") add "$2";;
  "rm") remove "$2";;
  "ls") list_save;;
  "ia") ignore_add "$2";;
  "irm") ignore_remove "$2";;
  "ils") list_ignore;;
  "s") sync;;
  "c") clean;;
  "h") help;;
  *) echo "Unknown mode, use h for list of modes";;
esac
