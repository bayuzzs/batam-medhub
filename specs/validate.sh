#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
tmp_dir="$(mktemp -d)"
trap 'rm -rf "$tmp_dir"' EXIT

cd "$repo_root"

npx --yes @redocly/cli lint specs/openapi.yaml specs/provider-openapi.yaml

for spec in specs/openapi.yaml specs/provider-openapi.yaml; do
  output="$tmp_dir/$(basename "$spec")"
  ruby -rjson -ryaml -e '
    source = File.expand_path(ARGV.fetch(0))
    output = ARGV.fetch(1)
    base = File.dirname(source)

    inline = lambda do |value|
      case value
      when Hash
        if value.keys == ["externalValue"] && value["externalValue"].start_with?("./")
          {"value" => JSON.parse(File.read(File.expand_path(value["externalValue"], base)))}
        else
          value.to_h { |key, child| [key, inline.call(child)] }
        end
      when Array
        value.map { |child| inline.call(child) }
      else
        value
      end
    end

    document = YAML.safe_load(File.read(source), aliases: true)
    File.write(output, YAML.dump(inline.call(document)))
  ' "$spec" "$output"

  npx --yes @redocly/cli lint "$output" --max-problems 500
done
