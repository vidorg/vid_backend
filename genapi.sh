python ./docs/parse_yaml.py main.go ./docs/api.yaml

python ./docs/to_html.py ./docs/api.yaml ./docs/api.html

# swagger-markdown -i ./docs/api.yaml -o ./docs/api.md