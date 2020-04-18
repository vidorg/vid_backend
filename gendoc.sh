echo Generating yaml and swagger...

apiparser.exe \
  --do_yaml \
  --main main.go \
  --dir . \
  --yaml_output ./docs/api.yaml \
  --ext go \
  --do_swag \
  --swag_output ./docs/api.html
#  --do_apib \
#  --apib_output ./docs/api.apib \
#  --need_content_type
