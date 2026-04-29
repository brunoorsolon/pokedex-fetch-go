#!/usr/bin/env bash
set -euo pipefail

# Prepare embedded data from a local pokemon-colorscripts checkout.
# Usage: POKEMON_COLORSCRIPTS=/path/to/pokemon-colorscripts ./scripts/prepare_data.sh

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SRC="${POKEMON_COLORSCRIPTS:-/workspace/pokemon-colorscripts-main}"

if [[ ! -f "$SRC/pokemon.json" || ! -d "$SRC/colorscripts" ]]; then
  echo "pokemon-colorscripts checkout not found at: $SRC" >&2
  exit 1
fi

mkdir -p "$ROOT/data/gen" \
  "$ROOT/data/sprites/small/regular" "$ROOT/data/sprites/small/shiny" \
  "$ROOT/data/sprites/large/regular" "$ROOT/data/sprites/large/shiny"

node - "$SRC/pokemon.json" "$ROOT/data" <<'NODE'
const fs = require('fs');
const [src, outDir] = process.argv.slice(2);
const input = JSON.parse(fs.readFileSync(src, 'utf8'));
const ranges = [[1,151],[152,251],[252,386],[387,493],[494,649],[650,721],[722,809],[810,905]];
function genFor(n) {
  for (let i = 0; i < ranges.length; i++) if (n >= ranges[i][0] && n <= ranges[i][1]) return i + 1;
  return 0;
}
const out = input.map((p, i) => ({
  name: p.name,
  number: i + 1,
  generation: genFor(i + 1),
  forms: p.forms || ['regular'],
  types: [],
  height: 0,
  weight: 0,
  stats: {hp: 0, attack: 0, defense: 0, sp_attack: 0, sp_defense: 0, speed: 0},
  is_legendary: false,
  is_mythical: false,
}));
fs.writeFileSync(`${outDir}/pokemon.json`, JSON.stringify(out, null, 2) + '\n');
for (let g = 1; g <= 8; g++) {
  fs.writeFileSync(`${outDir}/gen/gen${g}.txt`, out.filter(p => p.generation === g).map(p => p.name).join('\n') + '\n');
}
NODE

for size in small large; do
  for variant in regular shiny; do
    src_dir="$SRC/colorscripts/$size/$variant"
    dst_dir="$ROOT/data/sprites/$size/$variant"
    find "$dst_dir" -type f -name '*.gz' -delete
    while IFS= read -r -d '' file; do
      name="$(basename "$file")"
      gzip -c "$file" > "$dst_dir/$name.gz"
    done < <(find "$src_dir" -maxdepth 1 -type f -print0)
  done
done

echo "Prepared $(find "$ROOT/data" -type f | wc -l) embedded data files."
