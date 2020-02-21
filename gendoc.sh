echo Generating swagger yaml...
python3 ./docs/script/gen_yaml.py \
    -m ./main.go \
    -s . \
    -n true \
    -o ./docs/api.yaml \
    -e go

echo
echo Generating swagger html...
python3 ./docs/script/gen_swagger.py \
    -i ./docs/api.yaml \
    -o ./docs/api.html
